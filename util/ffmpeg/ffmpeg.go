package ffmpeg

import (
	"bufio"
	"bytes"
	"context"
	"embed"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/Carmen-Shannon/gopus"
)

//go:embed linux/ffmpeg
var linuxBinary embed.FS

//go:embed linux/ffmpeg-arm
var linuxArmBinary embed.FS

//go:embed windows/ffmpeg.exe
var windowsBinary embed.FS

//go:embed macos/ffmpeg
var macosBinary embed.FS

var (
	ffmpegPath string
	once       sync.Once
)

func getFFmpegPath() (string, error) {
	var err error
	once.Do(func() {
		ffmpegPath, err = extractBinary("ffmpeg")
		if err != nil {
			return
		}
	})
	return ffmpegPath, err
}

func extractBinary(name string) (string, error) {
	var fs embed.FS
	var path string
	switch runtime.GOOS {
	case "linux":
		switch runtime.GOARCH {
		case "arm", "arm64":
			path = fmt.Sprintf("linux/%s-arm", name)
			fs = linuxArmBinary
		default:
			path = fmt.Sprintf("linux/%s", name)
			fs = linuxBinary
		}
	case "windows":
		path = fmt.Sprintf("windows/%s.exe", name)
		fs = windowsBinary
	case "darwin":
		path = fmt.Sprintf("macos/%s", name)
		fs = macosBinary
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	tmpDir := os.TempDir()
	binaryPath := filepath.Join(tmpDir, filepath.Base(path))

	// Check if the binary already exists
	if _, err := os.Stat(binaryPath); err == nil {
		return binaryPath, nil
	}

	data, err := fs.ReadFile(path)
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(binaryPath, data, 0755); err != nil {
		return "", err
	}

	// un-embed the binaries, setting the embedded vars to nil
	macosBinary = embed.FS{}
	windowsBinary = embed.FS{}
	linuxBinary = embed.FS{}
	linuxArmBinary = embed.FS{}

	// make it purr
	data = nil
	runtime.GC()

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

// ConvertMp4ToMp3 will take a static file path ending in .mp4 and convert it into a temporary .mp3 file.
// The caller of this function should handle cleaning up the temporary file after it is done with it.
//
// Parameters:
//   - inputPath: the path to the mp4 file to convert.
//
// Returns:
//   - tempPath: the path to the temporary mp3 file.
//   - err: if an error occurs during the conversion process, this function will return an error.
func ConvertMp4ToMp3(inputPath string) (tempPath string, err error) {
	ffmpegPath, err := getFFmpegPath()
	if err != nil {
		return "", err
	}

	// Resolve the input path
	absInputPath, err := resolvePath(inputPath)
	if err != nil {
		return "", err
	}

	tempDir := os.TempDir()
	randomFileName := fmt.Sprintf("temp-%d.mp3", time.Now().UnixNano())
	tempFilePath := filepath.Join(tempDir, randomFileName)

	// Create the FFmpeg command
	cmd := exec.Command(
		ffmpegPath,
		"-i", absInputPath,
		"-vn",
		tempFilePath,
	)

	// Run the command and wait for it to complete
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to convert mp4 to mp3: %w", err)
	}

	// Return the path to the temporary file
	return tempFilePath, nil
}

// ConvertFileToPCM will take a static file path to an existing audio file in the appropriate format and convert it to PCM.
//
// This function is non-blocking and will send the PCM data to the output channel as long as the channel remains open, or the context isn't cancelled.
//
// Parameters:
//   - ctx: the context to listen for cancellation signals on. if the context is cancelled, this function will kill the ffmpeg process.
//   - inputPath: the path to the audio file to convert, valid audio formats are MP3, FLAC, WAV, and OGG. ultimately anything that ffmpeg can extract PCM from.
//   - outputChan: the channel you want to send the PCM data to. this is non-blocking and this channel will never be closed by this function.
//   - closeOutputChan: a function that will close the output channel when called. this is used to signal the end of the PCM data stream.
//
// Returns:
//   - <-chan struct{}: a channel that will be closed when the first packet is sent to the output channel.
//   - error: if an error occurs during the conversion process, this function will return an error. if the context is cancelled, this function will return nil.
func ConvertFileToPCM(ctx context.Context, inputPath string, outputChan chan []byte, closeOutputChan func()) (<-chan struct{}, error) {
	readySignal := make(chan struct{})
	closeReadySignal := sync.Once{}
	closeFunc := func() {
		closeReadySignal.Do(func() {
			close(readySignal)
		})
	}
	ffmpegPath, err := getFFmpegPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get ffmpeg path: %w", err)
	}

	// Resolve the input path
	absInputPath, err := resolvePath(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve input path: %w", err)
	}
	var tempPath string
	if strings.Contains(inputPath, ".mp4") {
		tempPath, err = ConvertMp4ToMp3(inputPath)
		if err != nil {
			return nil, fmt.Errorf("failed to convert mp4 to mp3: %w", err)
		}
		absInputPath = tempPath
	}

	ffmpegCmd := exec.Command(
		ffmpegPath,
		"-hide_banner",
		"-i", absInputPath,
		"-acodec", "pcm_s16le",
		"-f", "s16le",
		"-ar", "48000",
		"-ac", "2",
		"pipe:1",
	)

	pcmOut, err := ffmpegCmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	pcmBuf := bufio.NewReaderSize(pcmOut, 4096)
	// Start the ffmpeg command
	err = ffmpegCmd.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	go func() {
		defer ffmpegCmd.Process.Kill()
		defer ffmpegCmd.Process.Release()
		defer closeFunc()
		defer closeOutputChan()
		if tempPath != "" {
			defer os.Remove(tempPath)
		}

		for {
			select {
			case <-ctx.Done():
				return
			default:
				frameBuf := make([]byte, 3840)
				n, err := io.ReadFull(pcmBuf, frameBuf)
				if err != nil {
					if err == io.EOF || err == io.ErrUnexpectedEOF {
						if n > 0 {
							outputChan <- frameBuf[:n]
						}
						return
					}
					return
				}

				select {
				case <-ctx.Done():
					return
				case outputChan <- frameBuf[:n]:
					closeFunc()
				}
			}
		}
	}()

	return readySignal, nil
}

// ConvertPcmBytesToOpus will take bytes coming from the inputChan and convert them to Opus encoded audio frames. The Opus encoded frames are then sent through the outputChan.
// This function is non-blocking and will send the Opus data to the output channel as long as the channel remains open, or the context isn't cancelled.
// This function WILL close the `outputChan` given to it, when it detects that it reaches the end of the input channel, or the PCM data has been read entirely.
//
// Parameters:
//   - ctx: the context to listen for cancellation signals on. if the context is cancelled, this function will stop processing the input channel.
//   - inputChan: the channel you want to receive the PCM data from.
//   - outputChan: the channel you want to send the Opus data to. this channel will be closed by this function when the input channel is closed.
//   - closeOutputChan: a function that will close the output channel when called. this is used to signal the end of the Opus data stream.
//
// Returns:
//   - error: if an error occurs during the conversion process, this function will return an error. if the context is cancelled, this function will return nil.
func ConvertPcmBytesToOpus(ctx context.Context, inputChan chan []byte, outputChan chan []byte, closeOutputChan func()) error {
	sampleRate := 48000
	frameSize := 960
	channels := 2
	maxBytes := frameSize * channels * 2

	opusEncoder, err := gopus.NewEncoder(sampleRate, channels, gopus.Audio)
	if err != nil {
		return fmt.Errorf("failed to create Opus encoder: %w", err)
	}
	opusEncoder.SetBitrate(96000)
	opusEncoder.SetVbr(false)

	go func() {
		defer closeOutputChan()
		opusOutput := make([]byte, maxBytes)
		pcmBuf := make([]int16, frameSize*channels)
		for {
			select {
			case <-ctx.Done():
				return
			case frame, ok := <-inputChan:
				if !ok {
					return
				}
				frameReader := bytes.NewReader(frame)
				err = binary.Read(frameReader, binary.LittleEndian, &pcmBuf)
				if err != nil {
					if err == io.EOF || err == io.ErrUnexpectedEOF {
						return
					}
					return
				}

				opusOutput, err = opusEncoder.Encode(pcmBuf, frameSize, opusOutput)
				if err != nil {
					return
				}

				select {
				case <-ctx.Done():
					return
				case outputChan <- opusOutput:
				}
			}
		}
	}()
	return nil
}
