package client

import (
	"testing"
	"time"
)

func TestClient_EnqueueMessage(t *testing.T) {

}

func TestChannel(t *testing.T) {

	slow := make(chan int64, 2)
	stop := false

	fast := make(chan int64)

	done := make(chan struct{})

	out := make(chan int64)
	go func() {
		for i := range out {
			time.Sleep(time.Second * 1)
			t.Log("from", i)
		}
	}()

	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(time.Second * 1)
			t.Log("=========================================")
		}
	}()
	go func() {
		for i := 0; i < 60; i++ {
			time.Sleep(time.Millisecond * 500)
			if !stop {
				t.Log("======")
				slow <- int64(i)
			}
		}
	}()

	go func() {
		for i := 0; i < 4; i++ {
			time.Sleep(time.Millisecond * 500)
			t.Log("fast >>>")
			fast <- int64(i * -1)
		}
	}()

	go func() {
		time.Sleep(time.Second * 12)
		done <- struct{}{}
		stop = true
		close(slow)
	}()

	for stop := true; stop; {
		select {
		case m := <-slow:
			out <- m
		case m := <-fast:
			out <- m
		case <-done:
			t.Log("done")
			stop = false
		default:

		}
	}

	t.Log("Complete")
	time.Sleep(time.Second * 10)
}
