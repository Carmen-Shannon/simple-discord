package ffmpeg

import (
	"bufio"
	"bytes"
	"context"
	"embed"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"

	"github.com/Carmen-Shannon/gopus"
	"github.com/Carmen-Shannon/simple-discord/structs/voice"
)

//go:embed linux/ffmpeg linux/ffmpeg-arm windows/ffmpeg.exe macos/ffmpeg
//go:embed linux/ffprobe linux/ffprobe-arm windows/ffprobe.exe macos/ffprobe
var ffmpegFiles embed.FS

var (
	ffmpegPath  string
	ffprobePath string
	FfmpegCmd   *exec.Cmd
	FfprobeCmd  *exec.Cmd
	once        sync.Once
)

func getFFmpegPath() (string, error) {
	var err error
	once.Do(func() {
		ffmpegPath, err = extractBinary("ffmpeg")
		if err != nil {
			return
		}
		ffprobePath, err = extractBinary("ffprobe")
	})
	return ffmpegPath, err
}

func getFFprobePath() (string, error) {
	var err error
	once.Do(func() {
		ffmpegPath, err = extractBinary("ffmpeg")
		if err != nil {
			return
		}
		ffprobePath, err = extractBinary("ffprobe")
	})
	return ffprobePath, err
}

func extractBinary(name string) (string, error) {
	var path string
	switch runtime.GOOS {
	case "linux":
		switch runtime.GOARCH {
		case "arm", "arm64":
			path = fmt.Sprintf("linux/%s-arm", name)
		default:
			path = fmt.Sprintf("linux/%s", name)
		}
	case "windows":
		path = fmt.Sprintf("windows/%s.exe", name)
	case "darwin":
		path = fmt.Sprintf("macos/%s", name)
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	tmpDir := os.TempDir()
	binaryPath := filepath.Join(tmpDir, filepath.Base(path))

	// Check if the binary already exists
	if _, err := os.Stat(binaryPath); err == nil {
		return binaryPath, nil
	}

	data, err := ffmpegFiles.ReadFile(path)
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(binaryPath, data, 0755); err != nil {
		return "", err
	}

	return binaryPath, nil
}

func resolvePath(inputPath string) (string, error) {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Resolve the input path relative to the root project directory
	absPath := filepath.Join(cwd, inputPath)
	return absPath, nil
}

func ConvertFileToOpus(inputPath string, preserveMetadata bool, outputChan chan []byte, cancel context.CancelFunc) (*voice.AudioMetadata, error) {
	ffmpegPath, err := getFFmpegPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get ffmpeg path: %w", err)
	}

	// Resolve the input path
	absInputPath, err := resolvePath(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve input path: %w", err)
	}

	// Extract metadata from the original file
	var metadata *voice.AudioMetadata
	if preserveMetadata {
		metadata, err = GetOpusMetadataFromFile(absInputPath)
		if err != nil {
			return nil, fmt.Errorf("failed to get metadata: %w", err)
		}
	}

	// Create a shell command to run FFmpeg to convert MP3 to raw PCM
	FfmpegCmd = exec.Command(
		ffmpegPath,
		"-hide_banner",
		"-i", absInputPath,
		"-f", "s16le",
		"-ar", "48000",
		"-ac", "2",
		"pipe:1",
	)
	ffmpegOut, err := FfmpegCmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	ffmpegBuf := bufio.NewReaderSize(ffmpegOut, 16384)

	// Start the FFmpeg command
	err = FfmpegCmd.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	// Define frame size and channels
	frameSize := 960
	channels := 2
	maxBytes := frameSize * channels * 2

	// Create Opus encoder
	opusEncoder, err := gopus.NewEncoder(48000, 2, gopus.Audio)
	if err != nil {
		return nil, fmt.Errorf("failed to create Opus encoder: %w", err)
	}
	opusEncoder.SetBitrate(96000)
	opusEncoder.SetVbr(false)

	go func() {
		defer FfmpegCmd.Process.Kill()
		defer func() {
			if r := recover(); r != nil {
				if r == "send on closed channel" {
					return
				} else {
					fmt.Println(r)
				}
			}
		}()
		for {
			select {
			case <-outputChan:
				fmt.Println("CANCELLING FROM OUTPUT CHAN")
				return
			default:
				audiobuf := make([]int16, frameSize*channels)
				err = binary.Read(ffmpegBuf, binary.LittleEndian, &audiobuf)
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					close(outputChan)
					fmt.Println("CLOSING OUTPUT CHAN")
					break
				}
				if err != nil {
					fmt.Println("Error reading from ffmpeg stdout:", err)
					cancel()
					return
				}

				// Encode PCM to Opus
				out := make([]byte, maxBytes)
				out, err := opusEncoder.Encode(audiobuf, frameSize, out)
				if err != nil {
					fmt.Println("Error encoding PCM to Opus:", err)
					cancel()
					return
				}

				outputChan <- out
			}
		}
	}()

	return metadata, nil
}

func GetOpusMetadataFromBytes(input []byte) (*voice.AudioMetadata, error) {
	ffprobePath, err := getFFprobePath()
	if err != nil {
		return nil, fmt.Errorf("failed to get ffprobe path: %w", err)
	}

	ffmpegCmd := exec.Command(ffprobePath, "-v", "error", "-show_entries", "format=duration,bit_rate", "-show_entries", "stream=sample_rate,channels", "-of", "json", "pipe:0")
	ffmpegCmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	ffmpegCmd.Stdout = &out
	ffmpegCmd.Stderr = &stderr
	err = ffmpegCmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to get opus metadata: %w, output: %s", err, stderr.String())
	}

	var metadata struct {
		Format struct {
			Duration string `json:"duration"`
			BitRate  string `json:"bit_rate"`
		} `json:"format"`
		Streams []struct {
			SampleRate string `json:"sample_rate"`
			Channels   int    `json:"channels"`
		} `json:"streams"`
	}

	if err := json.Unmarshal(out.Bytes(), &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	duration, err := strconv.ParseFloat(metadata.Format.Duration, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse duration: %w", err)
	}

	bitrate, err := strconv.Atoi(metadata.Format.BitRate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bitrate: %w", err)
	}

	sampleRate, err := strconv.Atoi(metadata.Streams[0].SampleRate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sample rate: %w", err)
	}

	return &voice.AudioMetadata{
		DurationMs: duration * 1000, // Convert to milliseconds
		Bitrate:    bitrate,
		SampleRate: sampleRate,
		Channels:   metadata.Streams[0].Channels,
	}, nil
}

func GetOpusMetadataFromFile(filePath string) (*voice.AudioMetadata, error) {
	ffprobePath, err := getFFprobePath()
	if err != nil {
		return nil, fmt.Errorf("failed to get ffprobe path: %w", err)
	}

	FfprobeCmd = exec.Command(ffprobePath, "-v", "error", "-show_entries", "format=duration,bit_rate", "-show_entries", "stream=sample_rate,channels", "-of", "json", filePath)
	var out bytes.Buffer
	var stderr bytes.Buffer
	FfprobeCmd.Stdout = &out
	FfprobeCmd.Stderr = &stderr
	err = FfprobeCmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to get opus metadata: %w, output: %s", err, stderr.String())
	}

	var metadata struct {
		Format struct {
			Duration string `json:"duration"`
			BitRate  string `json:"bit_rate"`
		} `json:"format"`
		Streams []struct {
			SampleRate string `json:"sample_rate"`
			Channels   int    `json:"channels"`
		} `json:"streams"`
	}

	if err := json.Unmarshal(out.Bytes(), &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	duration, err := strconv.ParseFloat(metadata.Format.Duration, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse duration: %w", err)
	}

	bitrate, err := strconv.Atoi(metadata.Format.BitRate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bitrate: %w", err)
	}

	sampleRate, err := strconv.Atoi(metadata.Streams[0].SampleRate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sample rate: %w", err)
	}

	return &voice.AudioMetadata{
		DurationMs: duration * 1000, // Convert to milliseconds
		Bitrate:    bitrate,
		SampleRate: sampleRate,
		Channels:   metadata.Streams[0].Channels,
	}, nil
}
