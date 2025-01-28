package structs

import "sync"

type SyncGroup struct {
	CloseChannels map[string]*sync.Once
}

func NewSyncGroup() *SyncGroup {
	return &SyncGroup{
		CloseChannels: make(map[string]*sync.Once),
	}
}

func (sg *SyncGroup) AddChannel(name string) {
	sg.CloseChannels[name] = &sync.Once{}
}
