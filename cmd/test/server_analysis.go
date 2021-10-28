package test

import (
	"bytes"
	"fmt"
	"github.com/wcharczuk/go-chart"
	"go_im/im"
	"go_im/im/statistics"
	"go_im/pkg/db"
	"go_im/pkg/logger"
	"net/http"
	"os"
	"time"
)

var cTestPoints = 1

func RunAnalysisServer() {

	db.Init()

	done := make(chan struct{})

	go func() {
		defer func() {
			e := recover()
			if e != nil {

			}
		}()
		server := im.NewServer(im.Options{
			SvrType:       im.WebSocket,
			ApiImpl:       im.NewApiRouter(),
			ClientMgrImpl: im.NewClientManager(),
			GroupMgrImpl:  im.NewGroupManager(),
		})
		server.Serve("0.0.0.0", 8080)
	}()

	go func() {
		time.Sleep(time.Second * 3)
		http.Handle("/done", &doneHandler{done: done})
		http.Handle("/statistic", statistic{})
	}()

	go func() {
		tick := time.Tick(time.Second)
		tm := 0
		for range tick {
			tm++
			logger.D("%d", tm)
		}
	}()
	<-done
}

type statistic struct{}

func (s statistic) ServeHTTP(writer http.ResponseWriter, request *http.Request) {}

type doneHandler struct {
	done  chan struct{}
	times int
}

func (d *doneHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	d.times += 1
	if d.times >= cTestPoints {
		genStatisticsChart()
		d.done <- struct{}{}
	}
}

func genStatisticsChart() {
	st := statistics.GetStatistics()
	exportChart("input", "Msg Input", st.MsgInputMillSec)
	exportChart("output", "Msg Output", st.MsgOutPutMillSec)
	exportChart("count", "Msg I/O Count", st.MsgCountMillSec)
	exportChart("online", "Online", st.Online)
	exportChart("error", "Errors", st.ErrorsMillSec)
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
			Name:      "Time/Second",
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
	dir := "./analysis/" + now
	_ = os.MkdirAll(dir, os.ModePerm)

	f, err := os.Create(dir + "/" + title + ".png")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, _ = f.WriteAt(buffer.Bytes(), 0)
	_ = f.Close()
}
