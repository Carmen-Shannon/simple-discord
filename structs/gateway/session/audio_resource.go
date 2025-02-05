package session

import (
	"context"
	"fmt"
	"sync"

	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/util/ffmpeg"
)

type audioResource struct {
	mu     *sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	closeGroup structs.SyncGroup
	pcmStream  chan []byte
	opusStream chan []byte
}

type AudioResource interface {
	RegisterFile(path string)
	Exit()
	GetCtx() context.Context
	GetPcmStream() chan []byte
	GetOpusStream() chan []byte
	ClosePcmStream()
	CloseOpusStream()
}

func NewAudioResource() AudioResource {
	a := &audioResource{
		mu:         &sync.Mutex{},
		closeGroup: *structs.NewSyncGroup(),
		pcmStream:  make(chan []byte, 5),
		opusStream: make(chan []byte, 5),
	}
	a.ctx, a.cancel = context.WithCancel(context.Background())
	a.closeGroup.AddChannel("pcmStream")
	a.closeGroup.AddChannel("opusStream")
	return a
}

func (a *audioResource) RegisterFile(path string) {
	// Register the file to the PCM stream, start sending PCM data to the pcmStream
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		ready, err := ffmpeg.ConvertFileToPCM(a.ctx, path, a.pcmStream, a.ClosePcmStream)
		if err != nil {
			fmt.Printf("failed to register file: %v\n", err)
			a.Exit()
			return
		}

		select {
		case <-a.ctx.Done():
			wg.Done()
			return
		case <-ready:
			wg.Done()
		}
	}()

	go func() {
		wg.Wait()
		err := ffmpeg.ConvertPcmBytesToOpus(a.ctx, a.pcmStream, a.opusStream, a.CloseOpusStream)
		if err != nil {
			fmt.Printf("failed to convert pcm to opus: %v", err)
			return
		}
	}()
}

func (a *audioResource) Exit() {
	a.cancel()
	a.ClosePcmStream()
	a.CloseOpusStream()
}

func (a *audioResource) GetCtx() context.Context {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.ctx
}

func (a *audioResource) GetPcmStream() chan []byte {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.pcmStream
}

func (a *audioResource) GetOpusStream() chan []byte {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.opusStream
}

func (a *audioResource) ClosePcmStream() {
	a.closeGroup.CloseChannels["pcmStream"].Do(func() {
		close(a.pcmStream)
	})
}

func (a *audioResource) CloseOpusStream() {
	a.closeGroup.CloseChannels["opusStream"].Do(func() {
		close(a.opusStream)
	})
}
