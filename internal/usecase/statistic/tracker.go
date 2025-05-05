package statistic

import (
	"log/slog"
	"sync"
	"time"
)

const syncPeriod = time.Second * 30

type StorageInterface interface {
	AddStat(userId uint32, bytesIn, bytesOut uint64) error
}

type incomeStat struct {
	userId uint32
	in     uint64
	out    uint64
}

type Tracker struct {
	stats           sync.Map // map[uint32]*statistic
	consumeIncomeCh chan incomeStat
	consumeStopCh   chan struct{}
	syncStopCh      chan struct{}
	wg              sync.WaitGroup
	storage         StorageInterface
	logger          *slog.Logger
}

func New(storage StorageInterface, logger *slog.Logger) *Tracker {
	return &Tracker{
		consumeIncomeCh: make(chan incomeStat, 10000),
		consumeStopCh:   make(chan struct{}),
		syncStopCh:      make(chan struct{}),
		storage:         storage,
		logger:          logger,
	}
}

func (t *Tracker) Start() {
	t.wg.Add(2)

	go func() {
		defer t.wg.Done()

		t.consume()
	}()

	go func() {
		defer t.wg.Done()

		t.sync()
	}()
}

func (t *Tracker) Stop() {
	t.logger.Debug("Statistic tracker stopping")

	close(t.consumeStopCh)

	t.wg.Wait()

	t.commit()

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
			t.logger.Warn("Statistic tracker sync stopped")
			return
		}
	}
}

func (t *Tracker) cache(userId uint32, in, out uint64) {
	val, _ := t.stats.LoadOrStore(userId, &statistic{})
	st := val.(*statistic)
	st.Increment(in, out)
}

func (t *Tracker) commit() {
	t.stats.Range(func(key, value any) bool {
		userId := key.(uint32)
		stat := value.(*statistic)
		in, out := stat.GetAndClean()

		if in == 0 && out == 0 {
			t.stats.Delete(userId)
			return true
		}

		t.logger.Debug("Statistic tracker dumping", "user", userId, "in", in, "out", out)

		if err := t.storage.AddStat(userId, in, out); err != nil {
			t.logger.Debug("Statistic tracker dump", "err", err)
			stat.Increment(in, out)
		}

		return true
	})
}
