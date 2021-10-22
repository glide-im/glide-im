package im

import (
	"bytes"
	"github.com/wcharczuk/go-chart"
	"go_im/im/statistics"
	"go_im/pkg/db"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

//go test -v -run=TestServerPerf -memprofile=mem.out
func TestServerPerf(t *testing.T) {

	db.Init()

	server := NewServer(Options{
		SvrType:       WebSocket,
		ApiImpl:       NewApiRouter(),
		ClientMgrImpl: NewClientManager(),
		GroupMgrImpl:  NewGroupManager(),
	})

	go func() {
		server.Serve("0.0.0.0", 8080)
	}()

	go func() {
		tick := time.Tick(time.Second)
		tm := 0
		for range tick {
			tm++
			t.Log(tm)
		}
	}()

	var msgInputLine []int64
	var msgOutputLine []int64
	var msgCountLine []int64
	var onlineLine []int64

	// TODO fix incorrect avg msg compute
	go func() {
		tick := time.Tick(time.Millisecond * 100)
		for range tick {
			s := statistics.GetStatistics()
			onlineLine = append(onlineLine, s.ConnEnter.Get()-s.ConnExit.Get())

			var preInput int64 = 0
			if len(msgInputLine) > 0 {
				preInput = msgInputLine[len(msgInputLine)-1]
			}
			msgInputLine = append(msgInputLine, s.MessageInput.Get()-preInput)

			var preOutput int64 = 0
			if len(msgOutputLine) > 0 {
				preOutput = msgOutputLine[len(msgOutputLine)-1]
			}
			msgOutputLine = append(msgOutputLine, s.MessageOutput.Get()-preOutput)

			msgCountLine = append(msgCountLine, s.MessageOutput.Get()+s.MessageInput.Get())
		}
	}()

	done := make(chan struct{})
	handler := doneHandler{done: done}
	http.Handle("/done", handler)
	<-done

	exportChart("input", "Msg Input", msgInputLine)
	exportChart("output", "Msg Output", msgOutputLine)
	exportChart("count", "Msg I/O Count", msgCountLine)
	exportChart("online", "Online", onlineLine)
}

type doneHandler struct {
	done chan struct{}
}

func (d doneHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	d.done <- struct{}{}
}

func exportChart(title string, yName string, data []int64) {

	y := make([]float64, len(data))
	x := make([]float64, len(data))
	for i := 0; i < len(x); i++ {
		x[i] = float64(i)
		y[i] = float64(data[i])
	}

	graph := chart.Chart{
		Title: title,
		XAxis: chart.XAxis{
			Name:      "Time/100 MillSec",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		YAxis: chart.YAxis{
			Name:      yName,
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: x,
				YValues: y,
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
				},
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	_ = graph.Render(chart.PNG, buffer)

	_ = ioutil.WriteFile(title+".png", buffer.Bytes(), os.ModePerm)
}
