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

type Task struct {
	offset int
	s      *slot
	at     time.Time

	fn func()
	C  chan struct{}
}

func (s *Task) TTL() int64 {
	now := float64(time.Now().UnixNano())
	at := float64(s.at.UnixNano())
	return int64(math.Floor((at-now)/float64(time.Millisecond) + 1.0/2.0))
}

func (s *Task) call() {
	if s.s == nil {
		return
	}
	Executor(func() {
		s.Cancel()
		if s.fn != nil {
			s.fn()
		}
		s.C <- struct{}{}
	})
}

func (s *Task) Callback(f func()) {
	s.fn = f
}

func (s *Task) Cancel() {
	if s.s != nil {
		s.s.remove(s)
		s.s = nil
	}
}

type slot struct {
	index  int
	next   *slot
	len    int
	values map[*Task]interface{}

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
			values:    map[*Task]interface{}{},
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

func (s *slot) put(offset int, v *Task) int {
	if offset < 0 {
		panic("offset less the zero")
	}
	if !s.circulate && s.index == s.len && offset > 0 {
		return offset
	}
	if offset == 0 {
		s.m.Lock()
		s.values[v] = nil
		v.s = s
		s.m.Unlock()
		return 0
	}
	if offset >= s.len {
		return offset - s.len
	}
	return s.next.put(offset-1, v)
}

func (s *slot) isEmpty() bool {
	s.m.Lock()
	defer s.m.Unlock()
	return len(s.values) == 0
}

func (s *slot) callAndRm() {
	if s.isEmpty() {
		return
	}
	s.m.Lock()
	for k := range s.values {
		k.call()
	}
	s.m.Unlock()
}

func (s *slot) remove(v *Task) {
	s.m.Lock()
	delete(s.values, v)
	s.m.Unlock()
}

func (s *slot) valueArray() []*Task {
	var r []*Task
	s.m.Lock()
	for k := range s.values {
		r = append(r, k)
	}
	s.m.Unlock()
	return r
}

type wheel struct {
	slotCap int
	remain  int

	slot *slot

	parent *wheel
	child  *wheel
}

func newWheel(buckets int, dep int, wheels int, child *wheel) *wheel {
	wh := &wheel{
		slot:    newSlot(dep == wheels, buckets),
		slotCap: int(math.Pow(float64(buckets), float64(dep))) / buckets,
		child:   child,
	}
	if dep == wheels {
		wh.remain = wh.slotCap * buckets
	}
	if child != nil {
		child.parent = wh
	}
	return wh
}

func (w *wheel) tick() {
	if w.parent != nil {
		w.remain--
		if w.remain <= 0 {
			w.remain = w.slotCap * w.slot.len
		}
		w.parent.tick()
	}
}

func (w *wheel) move() bool {
	if w.child != nil {
		for _, v := range w.slot.valueArray() {
			w.slot.remove(v)
			w.child.put(v)
		}
		if w.child.move() {
			w.slot = w.slot.next
			for _, v := range w.slot.valueArray() {
				w.slot.remove(v)
				w.child.put(v)
			}
			return w.slot.index == 0
		} else {
			return false
		}
	} else {
		w.tick()
		w.slot = w.slot.next
		w.slot.callAndRm()
		return w.slot.index == 0
	}
}

func (w *wheel) put(v *Task) {

	s := int(math.Floor(float64(v.offset) / float64(w.slotCap)))
	if s == 0 {
		if w.child == nil {
			v.call()
		} else {
			w.child.put(v)
		}
	} else {
		if w.child != nil {
			v.offset = v.offset - ((s-1)*w.slotCap + w.child.remain - 1) - 1
		}
		w.slot.put(s, v)
	}
}

func (w *wheel) put2(v *Task) {

	s := int(math.Floor(float64(v.offset) / float64(w.slotCap)))
	sl := w.slot
	if s == 0 {
		if w.child != nil {
			if w.child.remain > v.offset {
				w.child.put2(v)
			} else {
				v.offset = v.offset - w.child.remain
				sl.put(1, v)
			}
		} else {
			sl.put(s, v)
			v.call()
		}
	} else {
		if w.child != nil {
			v.offset = v.offset - ((s-1)*w.slotCap + w.child.remain - 1) - 1
			if v.offset >= w.slotCap {
				s++
				v.offset = v.offset - w.slotCap
			}
		}
		sl.put(s, v)
	}
}

// TimingWheel the timing wheel ticker implementation
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
			w.parent = wh
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

func (w *TimingWheel) After(timeout time.Duration) *Task {
	if timeout >= w.maxTimeout {
		panic(fmt.Sprintf("maxTimeout=%d, current=%d", w.maxTimeout, timeout))
	}
	//offset := int(float64(TTL) / float64(w.interval))
	offset := int(math.Floor(float64(timeout.Milliseconds())/float64(w.interval.Milliseconds()) + 1.0/2.0))

	ch := make(chan struct{})

	t := &Task{
		offset: offset,
		C:      ch,
		at:     time.Now().Add(timeout),
	}
	w.wheel.put2(t)
	return t
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
