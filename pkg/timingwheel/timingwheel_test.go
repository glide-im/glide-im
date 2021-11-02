package timingwheel

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestNewTimingWheel(t *testing.T) {

	tw := NewTimingWheel(time.Millisecond*1000, 3, 20)
	tk := time.NewTicker(time.Second)

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
		sleepRndMilleSec(1000, 2000)
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				ch, v := tw.After(time.Second * time.Duration(10+rand.Int63n(30)))
				d := v.at.Unix() - time.Now().Unix()
				t.Log("after=", d)
				<-ch
				it := time.Now().Unix() - v.at.Unix()
				t.Log(d, it)
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
