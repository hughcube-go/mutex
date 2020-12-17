package mutex

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Mutex_Lock(t *testing.T) {
	a := assert.New(t)

	mu := NewMutex(1)

	go func() {
		mu.Lock()
		defer mu.Unlock()

		time.Sleep(2 * time.Second)
	}()

	time.Sleep(1 * time.Second)
	a.False(mu.TryLock())

	time.Sleep(2 * time.Second)
	a.True(mu.TryLock())
	defer mu.Unlock()
}

func Test_Mutex_TryLock(t *testing.T) {
	a := assert.New(t)

	mu := NewMutex(1)

	go func() {
		mu.Lock()
		defer mu.Unlock()

		time.Sleep(2 * time.Second)
	}()

	time.Sleep(1 * time.Second)
	a.False(mu.TryLock())

	time.Sleep(2 * time.Second)
	a.True(mu.TryLock())
	defer mu.Unlock()
}

func Test_Mutex_TryLockWithTimeout(t *testing.T) {
	a := assert.New(t)

	mu := NewMutex(1)

	go func() {
		mu.Lock()
		defer mu.Unlock()

		time.Sleep(2 * time.Second)
	}()

	time.Sleep(500 * time.Millisecond)
	timestamp := time.Now().UnixNano()
	a.False(mu.TryLockWithTimeout(1 * time.Second))
	a.True(int64(time.Second) < (time.Now().UnixNano() - timestamp))

	time.Sleep(2 * time.Second)
	a.True(mu.TryLock())
	defer mu.Unlock()
}
