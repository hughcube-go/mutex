package mutex

import (
	"sync"
)

type KeyMutex interface {
	Lock(key string)
	Unlock(key string)
}

type keyMutex struct {
	locks  map[string]*sync.RWMutex
	counts map[string]int64
	mu     *sync.Mutex
}

func NewKeyMutex() KeyMutex {
	km := &keyMutex{
		locks:  map[string]*sync.RWMutex{},
		counts: map[string]int64{},
		mu:     new(sync.Mutex),
	}
	return km
}

func (km *keyMutex) Lock(key string) {
	km.getMutexForLock(key).Lock()
}

func (km *keyMutex) Unlock(key string) {
	km.getMutexForUnlock(key).Unlock()
}

func (km *keyMutex) getMutexForLock(key string) *sync.RWMutex {
	km.mu.Lock()
	defer km.mu.Unlock()

	// get Mutex
	if _, ok := km.locks[key]; !ok {
		km.locks[key] = new(sync.RWMutex)
	}

	// incr count
	if _, ok := km.counts[key]; !ok {
		km.counts[key] = 0
	}
	km.counts[key]++

	return km.locks[key]
}

func (km *keyMutex) getMutexForUnlock(key string) *sync.RWMutex {
	km.mu.Lock()
	defer km.mu.Unlock()

	mu, ok := km.locks[key]
	if !ok {
		panic("no lock for " + key + " found")
	}

	km.counts[key]--
	if 0 >= km.counts[key] {
		delete(km.counts, key)
		delete(km.locks, key)
	}

	return mu
}
