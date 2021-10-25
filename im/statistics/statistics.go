package statistics

import (
	"github.com/panjf2000/ants/v2"
	"sync/atomic"
	"time"
)

type Statistics struct {
	messageInput  int64
	messageOutput int64
	errors        int64
	connEnter     int64
	connExit      int64

	StartUpAt        time.Time
	Online           []int64
	ErrorsMillSec    []int64
	MsgCountMillSec  []int64
	MsgInputMillSec  []int64
	MsgOutPutMillSec []int64
}

var statistics *Statistics
var pool *ants.Pool

func init() {
	statistics = &Statistics{
		messageInput:  0,
		messageOutput: 0,
		errors:        0,
		connEnter:     0,
		StartUpAt:     time.Now(),
	}
	pool, _ = ants.NewPool(100,
		ants.WithNonblocking(true),
		ants.WithPreAlloc(true),
	)

	go runStatistic()
}

func runStatistic() {
	duration := time.Millisecond * 100
	t := time.Tick(duration)
	for range t {
		input := atomic.SwapInt64(&statistics.messageInput, 0)
		output := atomic.SwapInt64(&statistics.messageOutput, 0)
		err := atomic.SwapInt64(&statistics.errors, 0)
		ol := atomic.LoadInt64(&statistics.connEnter) - atomic.LoadInt64(&statistics.connExit)

		statistics.ErrorsMillSec = append(statistics.ErrorsMillSec, err)
		statistics.Online = append(statistics.Online, ol)
		statistics.MsgInputMillSec = append(statistics.MsgInputMillSec, input)
		statistics.MsgOutPutMillSec = append(statistics.MsgOutPutMillSec, output)
		statistics.MsgCountMillSec = append(statistics.MsgCountMillSec, input+output)
	}
}

func SMsgInput() {
	incr(&statistics.messageInput)
}

func SMsgOutput() {
	incr(&statistics.messageOutput)
}

func SError(e error) {
	incr(&statistics.errors)
}

func SConnEnter() {
	incr(&statistics.connEnter)
}

func SConnExit() {
	incr(&statistics.connExit)
}

func GetStatistics() Statistics {
	return *statistics
}

func incr(s *int64) {
	_ = pool.Submit(func() {
		atomic.AddInt64(s, 1)
	})
}
