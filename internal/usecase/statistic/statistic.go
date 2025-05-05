package statistic

import (
	"sync"
)

type statistic struct {
	in  uint64
	out uint64
	mu  sync.Mutex
}

func (s *statistic) Increment(in, out uint64) {
	s.mu.Lock()
	s.in += in
	s.out += out
	s.mu.Unlock()
}

func (s *statistic) GetAndClean() (uint64, uint64) {
	s.mu.Lock()
	oldIn := s.in
	oldOut := s.out
	s.in = 0
	s.out = 0
	s.mu.Unlock()

	return oldIn, oldOut
}
