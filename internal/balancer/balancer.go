package balancer

import (
	"context"
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
	b.checkAll(context.Background())
	return b
}

func (b *Balancer) HealthCheck(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			b.checkAll(ctx)
		}
	}
}

func (b *Balancer) checkAll(ctx context.Context) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for _, backend := range b.backends {
		b.alive[backend] = isAlive(ctx, backend)
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

func isAlive(ctx context.Context, backend string) bool {
	d := &net.Dialer{Timeout: 2 * time.Second}
	conn, err := d.DialContext(ctx, "tcp", backend)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

func (b *Balancer) SetBackendAlive(backend string, alive bool) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.alive[backend] = alive
}

func (b *Balancer) GetAliveBackends() map[string]bool {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	result := make(map[string]bool)
	for k, v := range b.alive {
		result[k] = v
	}
	return result
}
