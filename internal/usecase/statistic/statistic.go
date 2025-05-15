package statistic

import (
	"sync/atomic"
)

type statistic struct {
	in  uint64
	out uint64
}

func (s *statistic) Increment(in, out uint64) {
	if in > 0 {
		atomic.AddUint64(&s.in, in)
	}

	if out > 0 {
		atomic.AddUint64(&s.out, out)
	}
}

func (s *statistic) GetAndClean() (in, out uint64) {
	in = atomic.SwapUint64(&s.in, 0)
	out = atomic.SwapUint64(&s.out, 0)

	return
}
