package session

import (
	"context"
	"sync"

	"github.com/Carmen-Shannon/simple-discord/util/ffmpeg"
)

type audioResource struct {
	mu     *sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	stream chan []byte
}

type AudioResource interface {
	RegisterFile(path string) error
	Exit()
	GetStream() chan []byte
}

func NewAudioResource() AudioResource {
	a := &audioResource{
		mu:     &sync.Mutex{},
		stream: make(chan []byte),
	}
	a.ctx, a.cancel = context.WithCancel(context.Background())
	return a
}

func (a *audioResource) RegisterFile(path string) error {
	err := ffmpeg.ConvertFileToOpus(path, a.stream, a.ctx)
	if err != nil {
		return err
	}
	return nil
}

func (a *audioResource) Exit() {
	a.cancel()
}

func (a *audioResource) GetStream() chan []byte {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.stream
}
