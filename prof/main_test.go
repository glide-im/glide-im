package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"go_im/im/api"
	"go_im/im/client"
	"go_im/im/dao"
	"go_im/im/dao/uid"
	"go_im/im/message"
	"go_im/pkg/db"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var uids []int64
var msgTo = map[int64]int64{}
var ucIds = map[int64]int64{}
var cids = map[int64]int64{}

var conns = map[int64]*websocket.Conn{}

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

func TestA(t *testing.T) {

	db.Init()
	initUsers(100)
	t.Log("uids=", uids)
	t.Log("user init complete")

	wgConn := sync.WaitGroup{}
	for _, uid_ := range uids {
		wgConn.Add(1)
		time.Sleep(time.Millisecond * 100)
		id := uid_
		go func() {
			serverIM(t, id)
			wgConn.Done()
		}()
	}
	wgConn.Wait()
	t.Log("connection establish complete, count:", len(conns))

	go func() {
		tick := time.Tick(time.Second)
		for range tick {
			t.Log("send:", atomic.LoadInt64(sendMsg), "-", "receive:", atomic.LoadInt64(receiveMsg))
		}
	}()

	wgMsg := sync.WaitGroup{}
	for i, conn := range conns {
		id := i
		c := conn
		wgMsg.Add(1)
		go func() {
			startMsg(id, 30, c)
			wgMsg.Done()
		}()
	}
	wgMsg.Wait()
	t.Log("msg send complete")

	time.Sleep(time.Second * 5)
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
	t.Log("done")
	time.Sleep(time.Second * 3)
	_, _ = http.Get("http://localhost:8080/done")
}

func serverIM(t *testing.T, uid int64) {

	con, _, err := dialer.Dial("ws://localhost:8080/ws", nil)

	if err != nil {
		t.Error(err)
		return
	}

	con.SetCloseHandler(func(code int, text string) error {
		t.Log("closed!", code, text)
		return nil
	})
	go func() {
		for !connClosed {
			_, _, e := con.ReadMessage()
			if e != nil && !connClosed {
				t.Log(e)
				if strings.Contains(e.Error(), "unexpected EOF") ||
					strings.Contains(e.Error(), "use of closed network connection") {
					break
				}
			} else {
				atomic.AddInt64(receiveMsg, 1)
			}
		}
	}()

	login := message.NewMessage(1, "api.test.login", api.TestLoginRequest{
		Uid:    uid,
		Device: 2,
	})
	_ = con.WriteJSON(login)
	conns[uid] = con
}

func startMsg(uid int64, count int, conn *websocket.Conn) {
	c := cids[uid]
	ucId := ucIds[uid]
	to := msgTo[uid]

	for i := 0; i < count; i++ {
		sleepRndMilleSec(200, 600)
		m := &client.SenderChatMessage{
			Cid:         c,
			UcId:        ucId,
			TargetId:    to,
			MessageType: 0,
			Message:     "hello-world",
			SendAt:      dao.Timestamp(time.Now()),
		}
		msg := message.NewMessage(0, message.ActionChatMessage, m)
		er := conn.WriteJSON(msg)
		if er != nil {
			fmt.Println(er.Error())
			break
		} else {
			atomic.AddInt64(sendMsg, 1)
		}
	}
}

func initUsers(userCount int) {

	for i := 0; i < userCount; i++ {
		uids = append(uids, uid.GenUid())
	}

	for i := 0; i < userCount; i++ {
		from := uids[i]
		to := uids[rand.Int63n(int64(userCount))]
		msgTo[from] = to

		chat, err := dao.ChatDao.CreateChat(dao.ChatTypeUser, from, to)
		if err != nil {
			panic(err)
		}
		cids[from] = chat.Cid

		userChat, err := dao.ChatDao.NewUserChat(chat.Cid, from, to, dao.ChatTypeUser)
		if err != nil {
			panic(err)
		}
		ucIds[from] = userChat.UcId
	}
}

func sleepRndMilleSec(start int32, end int32) {
	n := rand.Int31n(end - start)
	n = start + n
	time.Sleep(time.Duration(n) * time.Millisecond)
}
