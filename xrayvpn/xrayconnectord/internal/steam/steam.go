package steam

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	datasteam "github.com/realglebivanov/hstd/hstdlib/dataloader/steam"
)

type Steam struct {
	cidrs  atomic.Value
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func New() *Steam {
	s := &Steam{}
	s.cidrs.Store([]string{})
	return s
}

func (s *Steam) Get() []string {
	return s.cidrs.Load().([]string)
}

func (s *Steam) StartRefresh(interval time.Duration) {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	s.wg.Go(func() { s.refreshLoop(ctx, interval) })
}

func (s *Steam) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
	s.wg.Wait()
}

func (s *Steam) refreshLoop(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	s.load()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.load()
		}
	}
}

func (s *Steam) load() {
	cidrs, err := datasteam.Fetch()
	if err != nil {
		slog.Error("load Steam CIDRs", "err", err)
		return
	}
	s.cidrs.Store(cidrs)
	slog.Info("loaded Steam CIDRs", "count", len(cidrs))
}
