package timingwheel

import (
	"sync"
	"time"
)

type TimingWheel struct {
	m sync.Mutex

	interval   time.Duration
	ticker     *time.Ticker
	quit       chan struct{}
	maxTimeout time.Duration

	cs  []chan struct{}
	pos int
}

func NewTimingWheel(interval time.Duration, buckets int) *TimingWheel {
	w := new(TimingWheel)

	w.m = sync.Mutex{}
	w.interval = interval
	w.quit = make(chan struct{})
	w.pos = 0

	w.maxTimeout = interval * (time.Duration(buckets))

	w.cs = make([]chan struct{}, buckets)
	w.ticker = time.NewTicker(interval)

	for i := range w.cs {
		w.cs[i] = make(chan struct{})
	}
	go w.run()

	return w
}

func (w *TimingWheel) Stop() {
	close(w.quit)
}

func (w *TimingWheel) After(timeout time.Duration) <-chan struct{} {
	if timeout >= w.maxTimeout {
		panic("timeout too much, over max timeout")
	}

	index := int(timeout / w.interval)
	if 0 < index {
		index--
	}

	w.m.Lock()
	index = (w.pos + index) % len(w.cs)
	b := w.cs[index]
	w.m.Unlock()

	return b
}

func (w *TimingWheel) run() {
	for {
		select {
		case <-w.ticker.C:
			w.onTicker()
		case <-w.quit:
			w.ticker.Stop()
			return
		}
	}
}

func (w *TimingWheel) onTicker() {
	w.m.Lock()
	lastC := w.cs[w.pos]
	w.cs[w.pos] = make(chan struct{})
	w.pos = (w.pos + 1) % len(w.cs)
	w.m.Unlock()
	close(lastC)
}
