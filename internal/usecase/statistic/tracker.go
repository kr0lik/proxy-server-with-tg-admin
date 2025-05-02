package statistic

import (
	"log/slog"
	"sync"
	"time"
)

type StorageInterface interface {
	AddStat(userId uint32, bytesIn, bytesOut uint64) error
}

const syncPeriod = time.Second * 30

type stat struct {
	UserId uint32
	In     uint64
	Out    uint64
}

type Tracker struct {
	incomeCh chan stat
	stopCh   chan chan struct{}
	stats    sync.Map // map[uint32]*statistic
	storage  StorageInterface
	logger   *slog.Logger
}

func New(storage StorageInterface, logger *slog.Logger) *Tracker {
	return &Tracker{
		incomeCh: make(chan stat, 10000),
		stopCh:   make(chan chan struct{}),
		stats:    sync.Map{},
		storage:  storage,
		logger:   logger,
	}
}

func (t *Tracker) Start() {
	stopSyncCh := make(chan struct{})

	go func() {
		t.consume()
		close(stopSyncCh)
	}()

	go t.sync(stopSyncCh)

	select {
	case <-t.stopCh:
		t.logger.Debug("Statistic tracker stopping")
		close(t.incomeCh)
		return
	}
}

func (t *Tracker) Stop() {
	close(t.stopCh)
}

func (t *Tracker) Track(UserId uint32, In uint64, Out uint64) {
	t.incomeCh <- stat{UserId: UserId, In: In, Out: Out}
}

func (t *Tracker) consume() {
	t.logger.Debug("Statistic tracker consume starting")

	for {
		select {
		case dto, ok := <-t.incomeCh:
			if !ok {
				t.logger.Warn("Statistic tracker consume stopped")
				return
			}

			if dto.UserId == 0 {
				t.logger.Warn("Statistic tracker", "UserId", dto.UserId)
				continue
			}

			if dto.In == 0 && dto.Out == 0 {
				continue
			}

			t.cache(dto.UserId, dto.In, dto.Out)
		}
	}
}

func (t *Tracker) sync(stopSyncCh chan struct{}) {
	t.logger.Debug("Statistic tracker sync start")

	ticker := time.NewTicker(syncPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			t.commit()
		case <-stopSyncCh:
			t.logger.Warn("Statistic tracker sync stopped")
			t.commit()
			return
		}
	}
}

func (t *Tracker) cache(userId uint32, in, out uint64) {
	val, _ := t.stats.LoadOrStore(userId, &statistic{})
	st := val.(*statistic)
	st.Increment(in, out)
}

func (t *Tracker) delete(userId uint32) {
	val, ok := t.stats.Load(userId)
	if !ok {
		return
	}

	st := val.(*statistic)

	if st.in == 0 && st.out == 0 {
		t.stats.Delete(userId)

		t.logger.Debug("Statistic tracker cache", "delete", userId)
	}
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
