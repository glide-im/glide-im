package timingwheel

import (
	"fmt"
	"math"
	"sync"
	"time"
)

var Executor = func(f func()) {
	go f()
}

type value struct {
	offset int
	c      chan struct{}

	at time.Time
}

func (s *value) call() {
	Executor(func() {
		s.c <- struct{}{}
	})
}

type slot struct {
	index  int
	next   *slot
	len    int
	values map[*value]interface{}

	m         sync.Mutex
	circulate bool
}

func newSlot(circulate bool, len int) *slot {
	var head *slot
	var s *slot
	for i := 0; i < len; i++ {
		n := &slot{
			index:     i,
			len:       len,
			values:    map[*value]interface{}{},
			circulate: circulate,
			m:         sync.Mutex{},
		}
		if i == 0 {
			head = n
		} else {
			s.next = n
		}
		s = n
	}
	s.next = head
	return s
}

func (s *slot) put(offset int, v *value) int {
	if offset == 0 {
		s.m.Lock()
		s.values[v] = nil
		s.m.Unlock()
		return 0
	}
	if offset >= s.len {
		return offset - s.len
	}
	return s.next.put(offset-1, v)
}

func (s *slot) isEmpty() bool {
	return len(s.values) == 0
}

func (s *slot) callAndRm() {
	if s.isEmpty() {
		return
	}
	for k := range s.values {
		k.call()
		s.remove(k)
	}
}

func (s *slot) remove(v *value) {
	s.m.Lock()
	delete(s.values, v)
	s.m.Unlock()
}

func (s *slot) valueArray() []*value {
	var r []*value
	s.m.Lock()
	for k := range s.values {
		r = append(r, k)
	}
	s.m.Unlock()
	return r
}

type wheel struct {
	slotCap int

	slot  *slot
	child *wheel
}

func newWheel(buckets int, dep int, wheels int, child *wheel) *wheel {
	wh := &wheel{
		slot:    newSlot(dep == wheels, buckets),
		slotCap: int(math.Pow(float64(buckets), float64(dep))) / buckets,
		child:   child,
	}
	return wh
}

func (w *wheel) move() bool {
	if w.child != nil {
		if !w.slot.isEmpty() {
			for _, v := range w.slot.valueArray() {
				v.offset = v.offset % w.slotCap
				w.slot.remove(v)
				r := w.child.put(v)
				if r > 0 {
					v.offset = r
					w.slot.next.put(r, v)
				}
			}
		}
		if w.child.move() {
			w.slot = w.slot.next
			return w.slot.index == 0
		} else {
			return false
		}
	} else {
		w.slot.callAndRm()
		w.slot = w.slot.next
		return w.slot.index == 0
	}
}

func (w *wheel) put(v *value) int {
	offset := v.offset / w.slotCap
	if !w.slot.circulate {
		r := (w.slot.len - 1) - (w.slot.index + offset)
		if r < 0 {
			return -r
		}
	}
	if w.child == nil && offset == 0 {
		v.call()
	}
	return w.slot.put(offset, v)
}

type TimingWheel struct {
	interval   time.Duration
	ticker     *time.Ticker
	quit       chan struct{}
	maxTimeout time.Duration

	wheel *wheel
}

func NewTimingWheel(interval time.Duration, wheels int, slots int) *TimingWheel {
	tw := new(TimingWheel)

	tw.interval = interval
	tw.quit = make(chan struct{})
	s := int64(math.Pow(float64(wheels), float64(slots)))

	tw.maxTimeout = interval * time.Duration(s)
	tw.ticker = time.NewTicker(interval)

	var w *wheel
	for i := 1; i <= wheels; i++ {
		wh := newWheel(slots, i, wheels, nil)
		if w != nil {
			wh.child = w
		}
		w = wh
	}
	tw.wheel = w

	go tw.run()

	return tw
}

func (w *TimingWheel) Stop() {
	close(w.quit)
}

func (w *TimingWheel) After(timeout time.Duration) (<-chan struct{}, *value) {
	if timeout >= w.maxTimeout {
		panic(fmt.Sprintf("maxTimeout=%d, current=%d", w.maxTimeout, timeout))
	}
	index := int(timeout / w.interval)
	ch := make(chan struct{})
	s := &value{
		offset: index,
		c:      ch,
		at:     time.Now().Add(timeout),
	}
	w.wheel.put(s)
	return s.c, s
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
	w.wheel.move()
}
