package balancer

import (
	"sync/atomic"

	"github.com/Pablo-Loyola-Tantaruna/go-edge-gateway/pkg/models"
)

type ServerPool struct {
	Backends []*models.Backend
	Current  uint64
}

func (s *ServerPool) AddBackend(b *models.Backend) {
	s.Backends = append(s.Backends, b)
}

func (s *ServerPool) GetNextIndex() int {
	return int(atomic.AddUint64(&s.Current, uint64(1)) % uint64(len(s.Backends)))
}

func (s *ServerPool) GetNextPeerAfterFailure() *models.Backend {

	return s.GetNextPeer()
}

func (s *ServerPool) GetNextPeer() *models.Backend {
	next := s.GetNextIndex()
	l := len(s.Backends) + next

	for i := next; i < l; i++ {
		idx := i % len(s.Backends)
		if s.Backends[idx].IsAlive() {
			if i != next {
				atomic.StoreUint64(&s.Current, uint64(idx))
			}
			return s.Backends[idx]
		}
	}
	return nil
}
