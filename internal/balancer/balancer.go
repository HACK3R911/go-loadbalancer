package balancer

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type Balancer struct {
	backends []string
	alive    map[string]bool
	mutex    sync.RWMutex
	current  uint32
}

func New(backends []string) *Balancer {
	b := &Balancer{
		backends: backends,
		alive:    make(map[string]bool),
	}
	b.checkAll()
	return b
}

func (b *Balancer) HealthCheck(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		b.checkAll()
	}
}

func (b *Balancer) checkAll() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for _, backend := range b.backends {
		b.alive[backend] = isAlive(backend)
	}
}

func (b *Balancer) Next() (string, error) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if len(b.backends) == 0 {
		return "", errors.New("no backends configured")
	}

	start := int(atomic.LoadUint32(&b.current))
	for i := 0; i < len(b.backends); i++ {
		idx := (start + i) % len(b.backends)
		backend := b.backends[idx]
		if alive, ok := b.alive[backend]; ok && alive {
			atomic.StoreUint32(&b.current, uint32(idx+1))
			return backend, nil
		}
	}
	return "", errors.New("all backends are down")
}

func isAlive(backend string) bool {
	conn, err := net.DialTimeout("tcp", backend, 2*time.Second)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
