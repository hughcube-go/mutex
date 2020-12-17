package mutex

import (
	"context"
	"time"
)

type Mutex interface {
	Lock()

	TryLock() bool

	TryLockWithTimeout(time.Duration) bool

	TryLockWithContext(ctx context.Context) bool

	Unlock()
}

type chanMutex struct {
	lockChan chan struct{}
}

// NewChanMutex returns ChanMutex lock
func NewMutex(limit int) Mutex {
	return &chanMutex{
		lockChan: make(chan struct{}, limit),
	}
}

func (m *chanMutex) Lock() {
	m.lockChan <- struct{}{}
}

func (m *chanMutex) Unlock() {
	<-m.lockChan
}

func (m *chanMutex) TryLock() bool {
	select {
	case m.lockChan <- struct{}{}:
		return true
	default:
		return false
	}
}

func (m *chanMutex) TryLockWithContext(ctx context.Context) bool {
	select {
	case m.lockChan <- struct{}{}:
		return true
	case <-ctx.Done():
		return false
	}
}

func (m *chanMutex) TryLockWithTimeout(duration time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	return m.TryLockWithContext(ctx)
}
