package im

import (
	"bytes"
	"github.com/wcharczuk/go-chart"
	"go_im/im/statistics"
	"go_im/pkg/db"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

//go test -v -run=TestServerPerf -memprofile=mem.out
func TestServerPerf(t *testing.T) {

	db.Init()

	var closeAfter = time.Second * 180

	server := NewServer(Options{
		SvrType:       WebSocket,
		ApiImpl:       NewApiRouter(),
		ClientMgrImpl: NewClientManager(),
		GroupMgrImpl:  NewGroupManager(),
	})

	done := make(chan struct{})

	go func() {
		server.Serve("0.0.0.0", 8080)
	}()

	go func() {
		time.AfterFunc(closeAfter, func() {
			done <- struct{}{}
		})
	}()

	go func() {
		tick := time.Tick(time.Second)
		tm := 0
		for range tick {
			tm++
			t.Log("CountDown", 180-tm)
		}
	}()

	var msgInputLine []int64
	var msgOutputLine []int64
	var msgCountLine []int64
	var onlineLine []int64

	go func() {
		tick := time.Tick(time.Second)
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

	<-done

	exportChart("input", msgInputLine)
	exportChart("output", msgOutputLine)
	exportChart("count", msgCountLine)
	exportChart("online", onlineLine)
}

func exportChart(title string, data []int64) {

	y := make([]float64, len(data))
	x := make([]float64, len(data))
	for i := 0; i < len(x); i++ {
		x[i] = float64(i)
		y[i] = float64(data[i])
	}

	graph := chart.Chart{
		Title: title,
		XAxis: chart.XAxis{
			Name:      "sequence",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		YAxis: chart.YAxis{
			Name:      "value",
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
