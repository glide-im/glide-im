package main

import (
	"bytes"
	"fmt"
	"github.com/wcharczuk/go-chart"
	"go_im/im/dao/msgdao"
	"go_im/im/messaging"
	"go_im/im/statistics"
	"go_im/pkg/db"
	"go_im/pkg/logger"
	"net/http"
	"os"
	"time"
)

var cTestPoints = 1

func main() {
	RunAnalysisServer()
}

func RunAnalysisServer() {

	msgdao.MockChatMsg(time.Millisecond * 5)
	msgdao.MockCommDao()
	db.Init()
	messaging.Init()

	done := make(chan struct{})

	go func() {
		defer func() {
			e := recover()
			if e != nil {

			}
		}()
		RunTestServer()
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
			slog("%d", tm)
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
		slog("analysis image generated.")
		d.done <- struct{}{}
	}
}

func genStatisticsChart() {
	st := statistics.GetStatistics()
	//exportChart("input", "Msg Input", st.MsgInputMillSec)
	//exportChart("output", "Msg Output", st.MsgOutPutMillSec)
	//exportChart("count", "Msg I/O Count", st.MsgCountMillSec)
	//exportChart("online", "Online", st.Online)
	//exportChart("error", "Errors", st.ErrorsMillSec)

	exportChart(map[string][]int64{
		"input message":  st.MsgInputMillSec,
		"output message": st.MsgOutPutMillSec,
		"message count":  st.MsgCountMillSec,
		"online user":    st.Online,
		"error":          st.ErrorsMillSec,
	})
}

func exportChart(datas map[string][]int64) {

	var s []chart.Series

	colorIndex := 0
	for name, data := range datas {
		//transform time unit to second
		if true {
			var d []int64
			var ct int64 = 0
			for i, dat := range data {
				ct += dat
				if i%10 == 0 && i != 0 {
					if name == "online user" {
						ct = ct / 10
					}
					d = append(d, ct)
					ct = 0
				}
			}
			if ct != 0 {
				if name == "online user" {
					o := int64(len(data) % 10)
					if o == 0 {
						o = 1
					}
					ct = ct / o
				}
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

		color := chart.GetDefaultColor(colorIndex)
		colorIndex++
		s = append(s,
			chart.ContinuousSeries{
				Name:    name,
				XValues: x,
				YValues: y,
				Style: chart.Style{
					Show:        true,
					StrokeColor: color,
					FillColor:   color.WithAlpha(30),
				},
			},
		)
	}

	graph := chart.Chart{
		Width:  1440,
		Height: 900,
		Title:  "",
		XAxis: chart.XAxis{
			Name:      "Time/Second",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		YAxis: chart.YAxis{
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top:  20,
				Left: 20,
			},
		},
		Series: s,
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		fmt.Println(err.Error())
	}
	now := time.Now().Format("01-02_15_04_05")
	dir := "./analysis/" + now
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
	}
	f, err := os.Create(dir + "/" + "cps" + ".png")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = f.WriteAt(buffer.Bytes(), 0)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

}

func slog(format string, s ...interface{}) {
	if len(s) == 0 {
		logger.Zap.Sugar().Debugf(format)
	} else {
		logger.Zap.Sugar().Debugf(format, s...)
	}
}
