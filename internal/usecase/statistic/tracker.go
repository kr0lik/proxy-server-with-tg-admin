package statistic

import (
	"github.com/kagadar/go-syncmap"
	"log/slog"
	"sync"
	"time"
)

const syncPeriod = time.Second * 30
const consumeIncomeChBufferSize = 10000

type StorageInterface interface {
	AddStat(userId uint32, bytesIn, bytesOut uint64) error
}

type incomeStat struct {
	userId uint32
	in     uint64
	out    uint64
}

type Tracker struct {
	stats           syncmap.Map[uint32, *statistic]
	consumeIncomeCh chan incomeStat
	consumeStopCh   chan struct{}
	syncStopCh      chan struct{}
	wg              sync.WaitGroup
	storage         StorageInterface
	logger          *slog.Logger
}

func New(storage StorageInterface, logger *slog.Logger) *Tracker {
	return &Tracker{
		consumeIncomeCh: make(chan incomeStat, consumeIncomeChBufferSize),
		consumeStopCh:   make(chan struct{}),
		syncStopCh:      make(chan struct{}),
		storage:         storage,
		logger:          logger,
	}
}

func (t *Tracker) Start() {
	t.wg.Add(1)

	go func() {
		defer t.wg.Done()

		t.sync()
	}()

	t.wg.Add(1)

	go func() {
		defer t.wg.Done()

		t.consume()
	}()
}

func (t *Tracker) Stop() {
	t.logger.Debug("Statistic tracker stopping")

	close(t.consumeStopCh)

	t.wg.Wait()

	t.logger.Debug("Statistic tracker stopped")
}

func (t *Tracker) Track(userId uint32, in uint64, out uint64) {
	if userId == 0 {
		t.logger.Warn("Statistic tracker", "userId", userId)

		return
	}

	if in == 0 && out == 0 {
		return
	}

	select {
	case t.consumeIncomeCh <- incomeStat{userId: userId, in: in, out: out}:
	default:
		t.logger.Warn("Statistic tracker consumeIncomeCh buffer ends")
		t.cache(userId, in, out)
	}
}

func (t *Tracker) consume() {
	t.logger.Debug("Statistic tracker consume starting")

	for {
		select {
		case dto := <-t.consumeIncomeCh:
			t.cache(dto.userId, dto.in, dto.out)
		case <-t.consumeStopCh:
			for {
				select {
				case dto := <-t.consumeIncomeCh:
					t.cache(dto.userId, dto.in, dto.out)
				default:
					t.logger.Warn("Statistic tracker consume stopped")
					close(t.syncStopCh)

					return
				}
			}
		}
	}
}

func (t *Tracker) sync() {
	t.logger.Debug("Statistic tracker sync start")

	ticker := time.NewTicker(syncPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			t.commit()
		case <-t.syncStopCh:
			t.commit()
			t.logger.Warn("Statistic tracker sync stopped")

			return
		}
	}
}

func (t *Tracker) cache(userId uint32, in, out uint64) {
	st, _ := t.stats.LoadOrStore(userId, &statistic{})
	st.Increment(in, out)
}

func (t *Tracker) commit() {
	t.stats.Range(func(userId uint32, st *statistic) bool {
		in, out := st.GetAndClean()

		if in == 0 && out == 0 {
			t.stats.Delete(userId)

			return true
		}

		t.logger.Debug("Statistic tracker commit", "user", userId, "in", in, "out", out)

		if err := t.storage.AddStat(userId, in, out); err != nil {
			t.logger.Debug("Statistic tracker commit", "err", err)
			st.Increment(in, out)
		}

		return true
	})
}
