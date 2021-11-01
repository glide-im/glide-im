package timingwheel

import (
	"fmt"
	"go_im/pkg/logger"
	"math"
	"time"
)

type value struct {
	name   string
	offset int
	c      chan struct{}

	start int64
}

func (s *value) String() string {
	return fmt.Sprintf("value{name=%s,offset=%d,c=%values}", s.name, s.offset, s.c)
}

func (s *value) call() {
	s.c <- struct{}{}
}

type slot struct {
	index  int
	next   *slot
	len    int
	values map[*value]interface{}

	circulate bool
}

func newSlot(circulate bool, len int) *slot {
	slot1 := &slot{
		len:       len,
		values:    map[*value]interface{}{},
		circulate: circulate,
	}
	var s = slot1
	for i := 1; i < len; i++ {
		n := &slot{
			index:     i,
			len:       len,
			values:    map[*value]interface{}{},
			circulate: circulate,
		}
		s.next = n
		s = n
	}
	s.next = slot1
	return slot1
}

func (s *slot) put(offset int, v *value) int {
	if offset <= 0 {
		s.values[v] = nil
		return 0
	}
	if offset >= s.len {
		return offset - s.len
	}
	if !s.circulate && (s.index+offset) > s.len-1 {
		return (s.index + offset) - (s.len - 1)
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
		delete(s.values, k)
	}
}

func (s *slot) remove(v *value) {
	delete(s.values, v)
}

func (s *slot) valueArray() []*value {
	var r []*value
	for k := range s.values {
		r = append(r, k)
	}
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
		if w.child.move() {
			w.slot = w.slot.next
		}
		if !w.slot.isEmpty() {
			for _, v := range w.slot.valueArray() {
				v.offset = v.offset % w.slotCap
				w.slot.remove(v)
				w.child.put(v)
			}
		}
	} else {
		w.slot.callAndRm()
		w.slot = w.slot.next
	}
	if w.slot.next.index == 0 {
		w.slot = w.slot.next
		return true
	} else {
		return false
	}
}

func (w *wheel) put(v *value) {
	offset := v.offset / w.slotCap
	r := w.slot.put(offset, v)
	if r > 0 {
		logger.E("put failed", r, v)
	}
}

type TimingWheel struct {
	interval   time.Duration
	ticker     *time.Ticker
	quit       chan struct{}
	maxTimeout time.Duration

	wheel *wheel
}

func NewTimingWheel(interval time.Duration, buckets int) *TimingWheel {
	w := new(TimingWheel)

	w.interval = interval
	w.quit = make(chan struct{})
	w.maxTimeout = interval * (time.Duration(buckets ^ 3))
	w.ticker = time.NewTicker(interval)

	w1 := newWheel(buckets, 1, 3, nil)
	w2 := newWheel(buckets, 2, 3, w1)
	w.wheel = newWheel(buckets, 3, 3, w2)
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
	ch := make(chan struct{})
	s := &value{
		offset: index,
		c:      ch,
	}
	w.wheel.put(s)
	return nil
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
