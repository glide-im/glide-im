package timingwheel

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestNewTimingWheel(t *testing.T) {

	tw := NewTimingWheel(time.Second*1, 20)
	after := tw.After(time.Second * 3)
	tk := time.NewTicker(time.Second)
	c := 0

	status := func() {
		var w = tw.wheel
		for w != nil {
			t.Log(w.status())
			w = w.child
		}
	}
	for {
		select {
		case <-after:
			t.Log("done")
		case <-tk.C:
			c++
			t.Log("=>", c)
			status()
		}
	}
}

func (w *wheel) status() string {
	var s []string
	sl := w.slot
	for ; sl.index != sl.len-1; sl = sl.next {

	}
	for i := 0; i != sl.len; i++ {
		sl = sl.next
		if sl == w.slot {
			s = append(s, "*")
			continue
		}
		if sl.isEmpty() {
			s = append(s, "_")
		} else {
			s = append(s, "#")
		}
	}
	return fmt.Sprintf("value=%v", s)
}

func sleepRndMilleSec(start int32, end int32) {
	n := rand.Int31n(end - start)
	n = start + n
	time.Sleep(time.Duration(n) * time.Millisecond)
}
