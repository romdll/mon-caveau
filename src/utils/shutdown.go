package utils

import (
	"sync"
)

var (
	GracefulShutdownRegistry = NewGracefulShutdown()
)

type GracefulShutdown struct {
	mu       sync.Mutex
	channels map[string]chan struct{}
}

func NewGracefulShutdown() *GracefulShutdown {
	return &GracefulShutdown{
		channels: make(map[string]chan struct{}),
	}
}

func (gs *GracefulShutdown) Register(name string, ch chan struct{}) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.channels[name] = ch
}

func (gs *GracefulShutdown) TriggerShutdowns() {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	for name, ch := range gs.channels {
		close(ch)
		logger.Infof("Closed channel for: %s", name)
	}
}
