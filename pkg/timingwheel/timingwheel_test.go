package timingwheel

import (
	"testing"
	"time"
)

func TestTimingWheel_After(t *testing.T) {

	tw := NewTimingWheel(time.Second, 10)

	for range tw.After(time.Second * 3) {
		t.Log("#")
	}
	t.Log("done")
}
