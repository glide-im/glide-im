package statistics

import (
	"github.com/panjf2000/ants/v2"
	"go_im/im/comm"
	"time"
)

type Statistics struct {
	MessageInput  comm.AtomicInt64
	MessageOutput comm.AtomicInt64
	Errors        comm.AtomicInt64
	ConnEnter     comm.AtomicInt64
	ConnExit      comm.AtomicInt64
	StartUpAt     time.Time
}

var statistics *Statistics
var pool *ants.Pool

func init() {
	statistics = &Statistics{
		MessageInput:  0,
		MessageOutput: 0,
		Errors:        0,
		ConnEnter:     0,
		StartUpAt:     time.Now(),
	}
	pool, _ = ants.NewPool(100,
		ants.WithNonblocking(true),
		ants.WithPreAlloc(true),
	)
}

func SMsgInput() {
	incr(&statistics.MessageInput)
}

func SMsgOutput() {
	incr(&statistics.MessageOutput)
}

func SError(e error) {
	incr(&statistics.Errors)
}

func SConnEnter() {
	incr(&statistics.ConnEnter)
}

func SConnExit() {
	incr(&statistics.ConnExit)
}

func GetStatistics() Statistics {
	return *statistics
}

func incr(s *comm.AtomicInt64) {
	_ = pool.Submit(func() {
		s.Set(s.Get() + 1)
	})
}
