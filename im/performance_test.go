package im

import (
	"bytes"
	"fmt"
	"github.com/wcharczuk/go-chart"
	"go_im/im/statistics"
	"go_im/pkg/db"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestServerPerf(t *testing.T) {

	db.Init()

	done := make(chan struct{})
	http.Handle("/done", doneHandler{done: done})
	http.Handle("/statistic", statistic{})

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

	<-done
}

type statistic struct{}

func (s statistic) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	st := statistics.GetStatistics()
	fmt.Println(st.Online)
	exportChart("input", "Msg Input", st.MsgInputMillSec)
	exportChart("output", "Msg Output", st.MsgOutPutMillSec)
	exportChart("count", "Msg I/O Count", st.MsgCountMillSec)
	exportChart("online", "Online", st.Online)
	exportChart("error", "Errors", st.ErrorsMillSec)
}

type doneHandler struct {
	done chan struct{}
}

func (d doneHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	d.done <- struct{}{}
}

func TestE(t *testing.T) {
	exportChart("1", "1", []int64{1, 2, 3, 4, 3, 5, 6, 6, 0})
}
func exportChart(title string, yName string, data []int64) {

	// transform time unit to second
	if true {
		var d []int64
		var ct int64 = 0
		for i, dat := range data {
			ct += dat
			if i%10 == 0 && i != 0 {
				d = append(d, ct)
				ct = 0
			}
		}
		if ct != 0 {
			d = append(d, ct)
		}
		data = d
	}

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

	now := time.Now().Format("01-02_15_04_05")
	n := "./analysis/" + title + "_" + now + ".png"
	_ = os.MkdirAll(n, os.ModePerm)
	f, err := os.Create(n)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, _ = f.WriteAt(buffer.Bytes(), 0)
	_ = f.Close()
}
