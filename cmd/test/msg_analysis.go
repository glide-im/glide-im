package main

import (
	"bytes"
	"fmt"
	"github.com/wcharczuk/go-chart"
	"go_im/im"
	"go_im/im/statistics"
	"go_im/pkg/db"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main1() {
	//serveForAnalysis()
	exportChart("1", "1", []int64{1, 2, 3, 4, 3, 5, 6, 6, 0})
}

func serveForAnalysis() {

	db.Init()

	server := im.NewServer(im.Options{
		SvrType:       im.WebSocket,
		ApiImpl:       im.NewApiRouter(),
		ClientMgrImpl: im.NewClientManager(),
		GroupMgrImpl:  im.NewGroupManager(),
	})

	go func() {
		server.Serve("0.0.0.0", 8080)
	}()

	go func() {
		tick := time.Tick(time.Second)
		tm := 0
		for range tick {
			tm++
			fmt.Println("", tm)
		}
	}()

	done := make(chan struct{})
	handler := doneHandler{done: done}
	http.Handle("/done", handler)
	<-done

	s := statistics.GetStatistics()
	exportChart("input", "Msg Input", s.MsgInputMillSec)
	exportChart("output", "Msg Output", s.MsgOutPutMillSec)
	exportChart("count", "Msg I/O Count", s.MsgCountMillSec)
	exportChart("online", "Online", s.Online)
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

	e := ioutil.WriteFile(title+".png", buffer.Bytes(), os.ModePerm)
	fmt.Println(e)
}
