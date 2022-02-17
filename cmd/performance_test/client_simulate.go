package main

import (
	"bytes"
	"encoding/json"
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
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var uids []int64
var msgTo = map[int64]int64{}

var rLock = sync.RWMutex{}
var conns = map[int64]*websocket.Conn{}

var sentMessage = map[string]int64{}
var sentMsgMu = sync.Mutex{}

var messageDelay []int64

var sendMsg *int64
var receiveMsg *int64
var connClosed = false

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
var connToServerInterval = time.Millisecond * 10
var loginAfterConnInterval = time.Millisecond * 300
var totalUser = 10000
var msgPeerClient = 200
var msgSendIntervalFn = func() {
	sleepRndMilleSec(20, 100)
}

func RunClientMsg() {

	db.Init()
	//initUsers(totalUser)
	initUserNoDB(totalUser)

	slog("uids=%v", uids)
	slog("user init complete")

	wgConn := sync.WaitGroup{}
	for _, uid_ := range uids {
		wgConn.Add(1)
		time.Sleep(connToServerInterval)
		id := uid_
		go func() {
			connect(id)
			wgConn.Done()
		}()
	}
	wgConn.Wait()
	slog("connection establish complete, %d/%d", len(conns), totalUser)

	time.Sleep(time.Second * 1)
	go func() {
		tick := time.Tick(time.Second)
		count := totalUser * msgPeerClient
		for range tick {
			sent := atomic.LoadInt64(sendMsg)
			p := float64(sent) / float64(count)
			slog("Send:%d Receive:%d  Progress:%d%%", sent, atomic.LoadInt64(receiveMsg), int64(p*100))
		}
	}()

	wgMsg := sync.WaitGroup{}
	for i, conn := range conns {
		id := i
		c := conn
		wgMsg.Add(1)
		go func() {
			startMsg(id, msgPeerClient, c)
			wgMsg.Done()
		}()
	}
	wgMsg.Wait()
	slog("msg send complete")

	time.Sleep(time.Second * 1)
	connClosed = true
	for _, conn := range conns {
		c := conn
		go func() {
			m := message.NewMessage(0, "api.test.signout", "")
			_ = c.WriteJSON(m)
			time.Sleep(time.Millisecond * 100)
			_ = c.Close()
		}()
	}
	slog("done")
	_, _ = http.Get("http://" + host + ":8080/statistic")
	time.Sleep(time.Second * 3)
	_, _ = http.Get("http://" + host + ":8080/done")

	//exportMessageDelayChart()
}

func connect(uid int64) {

	con, _, err := dialer.Dial("ws://"+host+":8080/ws", nil)

	if err != nil {
		slog("err:%v", err.Error())
		return
	}

	con.SetCloseHandler(func(code int, text string) error {
		slog("err:%v", text)
		return nil
	})
	go func() {
		for !connClosed {
			_, b, e := con.ReadMessage()
			if e != nil && !connClosed {
				logger.W(e.Error())
				break
			} else {
				atomic.AddInt64(receiveMsg, 1)
				handleReceivedMessage(b)
			}
		}
	}()

	time.Sleep(loginAfterConnInterval)
	login := message.NewMessage(1, "api.test.login", test.TestLoginRequest{
		Uid:    uid,
		Device: 2,
	})
	_ = con.WriteJSON(login)
	rLock.Lock()
	conns[uid] = con
	rLock.Unlock()
}

type Message struct {
	Action string
	Data   string
}

func handleReceivedMessage(b []byte) {
	recvTime := time.Now()
	go func() {
		m := Message{}
		er := json.Unmarshal(b, &m)
		if er != nil {
			return
		}
		m2 := &message.DownChatMessage{}
		er = json.Unmarshal([]byte(m.Data), m2)
		if er != nil {
			return
		}
		sentMsgMu.Lock()
		v, ok := sentMessage[m2.Content]
		sentMsgMu.Unlock()
		if ok {
			messageDelay = append(messageDelay, recvTime.UnixNano()-v)
			sentMsgMu.Lock()
			delete(sentMessage, m2.Content)
			sentMsgMu.Unlock()
		}
	}()
}

func startMsg(uid int64, count int, conn *websocket.Conn) {
	to := msgTo[uid]

	for i := 0; i < count; i++ {
		msgSendIntervalFn()
		u := strconv.FormatInt(time.Now().UnixNano(), 10) + genRndString(10)
		m := &message.UpChatMessage{
			Mid:     1,
			CSeq:    1,
			To:      to,
			Type:    1,
			Content: u,
			CTime:   time.Now().Unix(),
		}
		msg := message.NewMessage(0, message.ActionChatMessage, m)
		s, _ := msg.Serialize()
		er := conn.WriteMessage(websocket.TextMessage, s)
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

func initUsers(userCount int) {

	for i := 0; i < userCount; i++ {
		uids = append(uids, uid.GenUid())
	}
	wg := sync.WaitGroup{}
	m1 := sync.Mutex{}

	for i := 0; i < userCount; i++ {
		wg.Add(1)
		i2 := i
		go func() {
			from := uids[i2]
			to := uids[rand.Int63n(int64(userCount))]
			m1.Lock()
			msgTo[from] = to
			m1.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
}

func initUserNoDB(count int) {
	for i := 0; i < count; i++ {
		uids = append(uids, uid.GenUid())
	}
	for i := 0; i < count; i++ {
		from := uids[i]
		to := uids[rand.Int63n(int64(count))]
		msgTo[from] = to
	}
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
	_, _ = f.WriteAt(buffer.Bytes(), 0)
	_ = f.Close()
}
