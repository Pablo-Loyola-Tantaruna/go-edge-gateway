package balancer

import (
	"sync"
	"sync/atomic"

	"github.com/Pablo-Loyola-Tantaruna/go-edge-gateway/pkg/models"
)

type ServerPool struct {
	Backends []*models.Backend
	mu       sync.RWMutex
	Current  uint64
}

func (s *ServerPool) AddBackend(b *models.Backend) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Backends = append(s.Backends, b)
}

func (s *ServerPool) getNextIndex() int {
	return int(atomic.AddUint64(&s.Current, uint64(1)) % uint64(len(s.Backends)))
}

func (s *ServerPool) GetNextPeerAfterFailure() *models.Backend {

	return s.GetNextPeer()
}

func (s *ServerPool) GetNextPeer() *models.Backend {
	s.mu.RLock()
	defer s.mu.RUnlock()

	n := len(s.Backends)
	if n == 0 {
		return nil
	}

	next := s.getNextIndex()
	l := n + next

	for i := next; i < l; i++ {
		idx := i % n
		if s.Backends[idx].IsAlive() {
			if i != next {
				atomic.StoreUint64(&s.Current, uint64(idx))
			}
			return s.Backends[idx]
		}
	}
	return nil
}

func (s *ServerPool) UpdateBackends(backends []*models.Backend) {
	s.mu.Lock()
	s.Backends = backends
	atomic.StoreUint64(&s.Current, 0)
	s.mu.Unlock()
}
