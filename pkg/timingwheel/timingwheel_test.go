package timingwheel

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestNewTimingWheel(t *testing.T) {

	tw := NewTimingWheel(time.Millisecond*100, 3, 20)
	runAt := time.Now()
	tk := time.NewTicker(time.Millisecond * 1000)

	status := func() {
		var w = tw.wheel
		for w != nil {
			//t.Log(w.status())
			w = w.child
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		//sleepRndMilleSec(1000, 2000)
		time.Sleep(time.Second * 3)
		for i := 0; i < 10000; i++ {
			wg.Add(1)
			go func() {
				sleepRndMilleSec(1000, 10000)
				tas := tw.After(time.Second * time.Duration(rand.Int63n(20)))
				after := tas.at.Unix() - time.Now().Unix()
				addAt := (time.Now().UnixNano() - runAt.UnixNano()) / int64(time.Millisecond)
				<-tas.C
				ttl := tas.TTL()
				if ttl > 100 {
					//t.Log("addAt=", addAt, "after=", after, "error=", err, "run=", time.Now().Unix()-runAt.Unix())
					t.Log("after:", after, "error:", ttl, addAt)
				} else {
					//t.Log("after:", after, "error:", "0")
				}
				wg.Done()
			}()
		}
		wg.Done()
	}()

	go func() {
		for {
			select {
			case <-tk.C:
				//t.Log("---------------------------------------------------------------------------")
				status()
			}
		}
	}()
	wg.Wait()
}

func (s *slot) tasks() [][]*Task {
	sl := s
	for sl.index != 0 {
		sl = sl.next
	}
	var t [][]*Task

	t = append(t, sl.valueArray())
	sl = sl.next
	if sl.index != 0 {
		t = append(t, sl.valueArray())
		sl = sl.next
	}
	t = append(t, sl.valueArray())
	return t
}

func (w *wheel) status() string {
	var s []string
	sl := w.slot
	for ; sl.index != sl.len-1; sl = sl.next {

	}
	for i := 0; i != sl.len; i++ {
		sl = sl.next
		if sl.index == w.slot.index {
			s = append(s, strconv.Itoa(i))
			continue
		}
		if sl.isEmpty() {
			s = append(s, "_")
		} else {
			s = append(s, "#")
		}
	}

	var ts []string
	for _, tasks := range w.slot.tasks() {
		var tt []string
		for _, t := range tasks {
			tt = append(tt, strconv.Itoa(t.offset))
		}
		ts = append(ts, fmt.Sprintf("%v", tt))
	}

	return fmt.Sprintf("%v %v %d", s, ts, w.remain)
}

func sleepRndMilleSec(start int32, end int32) {
	n := rand.Int31n(end - start)
	n = start + n
	time.Sleep(time.Duration(n) * time.Millisecond)
}
