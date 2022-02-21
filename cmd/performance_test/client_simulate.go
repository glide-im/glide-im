package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/wcharczuk/go-chart"
	"go_im/im/api/test"
	"go_im/im/dao/uid"
	"go_im/im/message"
	"go_im/pkg/db"
	"go_im/pkg/logger"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var uids []int64

var rLock = sync.RWMutex{}
var conns = map[int64]*websocket.Conn{}

var sentMessage = map[string]int64{}
var sentMsgMu = sync.RWMutex{}

var messageDelay []int64

var sendMsg *int64
var receiveMsg *int64

func init() {
	var s int64 = 0
	sendMsg = &s
	var r int64 = 0
	receiveMsg = &r
}

var dialer = websocket.Dialer{
	HandshakeTimeout: time.Minute,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
}

//var host = "192.168.1.162"

var host = "127.0.0.1"
var onlinePeerSecond = time.Millisecond * 20
var loginAfterConnected = time.Millisecond * 100
var totalUser = 10
var msgPeerClient = 400

func RunClientMsg() {

	db.Init()
	for i := 0; i < totalUser; i++ {
		uids = append(uids, uid.GenUid())
	}
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	slog("uids=%v", uids)
	slog("uid init complete, start simulation")

	go func() {
		tick := time.Tick(time.Second)
		count := totalUser * msgPeerClient
		for range tick {
			time.Unix(1, 1)
			sent := atomic.LoadInt64(sendMsg)
			progress := float64(sent) / float64(count)
			conn := len(conns)
			slog("Conn:%d\tSend:%d\tReceive:%d\tProgress:%d%%", conn, sent, atomic.LoadInt64(receiveMsg), int64(progress*100))
		}
	}()

	wgConn := sync.WaitGroup{}
	for _, uid_ := range uids {
		wgConn.Add(1)
		time.Sleep(onlinePeerSecond)
		id := uid_
		go func() {
			closed := make(chan interface{})
			conn := connect(id, closed)
			rLock.Lock()
			conns[id] = conn
			rLock.Unlock()
			//defer func() {
			//	e := recover()
			//	if e != nil {
			//		slog("err:%v", e)
			//	}
			//	wgConn.Done()
			//}()
			startMsg(msgPeerClient, conn, 100, 50)
			//rLock.RLock()
			//delete(conns, id)
			//rLock.RUnlock()
			exit(id, conn)
			closed <- struct{}{}
			wgConn.Done()
		}()
	}
	wgConn.Wait()

	slog("msg send complete")
	_, _ = http.Get("http://" + host + ":8080/statistic")
	time.Sleep(time.Second * 3)
	_, _ = http.Get("http://" + host + ":8080/done")
	exportMessageDelayChart()
}

func exit(id int64, conn *websocket.Conn) {
	if conn == nil {
		rLock.RLock()
		conn = conns[id]
		rLock.RUnlock()
	}
	if conn == nil {
		return
	}
	m := message.NewMessage(0, "api.test.signout", "")
	_ = conn.WriteJSON(m)
	time.Sleep(time.Millisecond * 100)
	_ = conn.Close()
}

func connect(uid int64, closed chan interface{}) *websocket.Conn {

	con, _, err := dialer.Dial("ws://"+host+":8080/ws", nil)

	if err != nil {
		slog("err:%v", err.Error())
		return nil
	}

	con.SetCloseHandler(func(code int, text string) error {
		slog("err:%v", text)
		return nil
	})
	go func() {
		for {
			select {
			case <-closed:
				break
			default:
				_, b, e := con.ReadMessage()
				if e != nil {
					v, ok := <-closed
					if v != nil && ok {
						break
					}
					logger.W(e.Error())
					break
				} else {
					atomic.AddInt64(receiveMsg, 1)
					handleReceivedMessage(b)
				}
			}
		}
	}()

	time.Sleep(loginAfterConnected)
	login := message.NewMessage(1, "api.test.login", test.TestLoginRequest{
		Uid:    uid,
		Device: 2,
	})
	_ = con.WriteJSON(login)
	return con
}

func handleReceivedMessage(b []byte) {
	recvTime := time.Now()
	go func() {
		m := message.Message{}
		er := c.Decode(b, &m)
		if er != nil {
			return
		}
		m2 := &message.DownChatMessage{}
		er = m.DeserializeData(m2)
		if er != nil {
			return
		}
		sentMsgMu.RLock()
		v, ok := sentMessage[m2.Content]
		sentMsgMu.RUnlock()
		if ok {
			messageDelay = append(messageDelay, recvTime.UnixNano()-v)
			sentMsgMu.Lock()
			delete(sentMessage, m2.Content)
			sentMsgMu.Unlock()
		}
	}()
}

var c = message.JsonCodec{}

func startMsg(count int, conn *websocket.Conn, end, start int32) {

	for i := 0; i < count; i++ {
		//msgSendIntervalFn()
		n := rand.Int31n(end - start)
		n = start + n
		time.Sleep(time.Duration(n) * time.Millisecond)

		rLock.RLock()
		var to int64 = 0
		for k := range conns {
			to = k
			break
		}
		if to == 0 {
			panic("conn is empty")
		}
		rLock.RUnlock()
		u := strconv.FormatInt(time.Now().UnixNano(), 10) + genRndString(10)
		m := &message.UpChatMessage{
			Mid:     1,
			Seq:     1,
			To:      to,
			Type:    1,
			Content: u,
			SendAt:  time.Now().Unix(),
		}
		msg := message.NewMessage(0, message.ActionChatMessage, m)
		jsonMsg, err := c.Encode(msg)
		if err != nil {
			panic(err)
		}
		er := conn.WriteMessage(websocket.TextMessage, jsonMsg)
		go func() {
			sentMsgMu.Lock()
			sentMessage[u] = time.Now().UnixNano()
			sentMsgMu.Unlock()
		}()
		if er != nil {
			logger.E(er.Error())
			break
		} else {
			atomic.AddInt64(sendMsg, 1)
		}
	}
}

var (
	table = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func genRndString(length int) string {
	res := ""
	for i := 0; i < length; i++ {
		idx := rand.Intn(62)
		res = res + table[idx:idx+1]
	}
	return res
}

func sleepRndMilleSec(start int32, end int32) {
	n := rand.Int31n(end - start)
	n = start + n
	time.Sleep(time.Duration(n) * time.Millisecond)
}

func exportMessageDelayChart() {

	var v []chart.Value

	delays := []int64{20, 50, 100, 200, 300, 500, 800, 1000, 1500, 2000, 3000, 4000, -1}
	var counts = make([]float64, len(delays))

	for _, delay := range messageDelay {
		var index = len(delays) - 1
		for idx, d := range delays {
			if delay <= int64(time.Millisecond)*d {
				index = idx
				break
			}
		}
		counts[index] = counts[index] + 1
	}

	for i, val := range counts {
		label := strconv.FormatInt(delays[i], 10)
		v = append(v, chart.Value{
			Label: label + "ms",
			Value: val,
		})
	}

	graph := chart.BarChart{
		Title:      "Message Delay Distribution",
		TitleStyle: chart.StyleShow(),
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Height:   1440,
		BarWidth: 900,
		XAxis: chart.Style{
			Show: true,
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		Bars: v,
	}

	buffer := bytes.NewBuffer([]byte{})
	_ = graph.Render(chart.PNG, buffer)

	now := time.Now().Format("01-02_15_04_05")
	dir := "./analysis/" + now
	_ = os.MkdirAll(dir, os.ModePerm)

	f, err := os.Create(dir + "/" + "msg_delay" + ".png")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, _ = f.Write(buffer.Bytes())
	_ = f.Close()
}
