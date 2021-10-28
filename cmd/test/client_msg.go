package test

import (
	"github.com/gorilla/websocket"
	"go_im/im/api"
	"go_im/im/dao"
	"go_im/im/dao/uid"
	"go_im/im/message"
	"go_im/pkg/db"
	"go_im/pkg/logger"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var uids []int64
var msgTo = map[int64]int64{}
var ucIds = map[int64]int64{}
var cids = map[int64]int64{}

var rLock = sync.RWMutex{}
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

var host = "192.168.1.123"

func RunClientMsg() {

	db.Init()
	userCount := 2000

	//initUsers(userCount)
	initUserNoDB(userCount)

	logger.D("uids=%v", uids)
	logger.D("user init complete")

	wgConn := sync.WaitGroup{}
	for _, uid_ := range uids {
		wgConn.Add(1)
		time.Sleep(time.Millisecond * 10)
		id := uid_
		go func() {
			serverIM(id)
			wgConn.Done()
		}()
	}
	wgConn.Wait()
	logger.D("connection establish complete, %d/%d", len(conns), userCount)

	time.Sleep(time.Second * 1)
	go func() {
		tick := time.Tick(time.Second)
		for range tick {
			logger.D("Send:%d  Receive:%d", atomic.LoadInt64(sendMsg), atomic.LoadInt64(receiveMsg))
		}
	}()

	wgMsg := sync.WaitGroup{}
	for i, conn := range conns {
		id := i
		c := conn
		wgMsg.Add(1)
		go func() {
			startMsg(id, 50, c)
			wgMsg.Done()
		}()
	}
	wgMsg.Wait()
	logger.D("msg send complete")

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
	logger.D("done")
	_, _ = http.Get("http://" + host + ":8080/statistic")
	time.Sleep(time.Second * 3)
	_, _ = http.Get("http://" + host + ":8080/done")
}

func serverIM(uid int64) {

	con, _, err := dialer.Dial("ws://"+host+":8080/ws", nil)

	if err != nil {
		logger.W(err.Error())
		return
	}

	con.SetCloseHandler(func(code int, text string) error {
		logger.W(text)
		return nil
	})
	go func() {
		for !connClosed {
			_, _, e := con.ReadMessage()
			if e != nil && !connClosed {
				logger.W(e.Error())
				break
			} else {
				atomic.AddInt64(receiveMsg, 1)
			}
		}
	}()

	time.Sleep(time.Millisecond * 300)
	login := message.NewMessage(1, "api.test.login", api.TestLoginRequest{
		Uid:    uid,
		Device: 2,
	})
	_ = con.WriteJSON(login)
	rLock.Lock()
	conns[uid] = con
	rLock.Unlock()
}

func startMsg(uid int64, count int, conn *websocket.Conn) {
	c := cids[uid]
	ucId := ucIds[uid]
	to := msgTo[uid]

	for i := 0; i < count; i++ {
		sleepRndMilleSec(60, 100)
		m := &message.SenderChatMessage{
			Cid:         c,
			UcId:        ucId,
			TargetId:    to,
			MessageType: 0,
			Message:     " hello-world hello-world hello-world hello-world hello-world",
			SendAt:      time.Now().Unix(),
		}
		msg := message.NewMessage(0, message.ActionChatMessage, m)
		s, _ := msg.Serialize()
		er := conn.WriteMessage(websocket.TextMessage, s)
		if er != nil {
			logger.E(er.Error())
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
	wg := sync.WaitGroup{}
	m1 := sync.Mutex{}
	m2 := sync.Mutex{}
	m3 := sync.Mutex{}

	for i := 0; i < userCount; i++ {
		wg.Add(1)
		i2 := i
		go func() {
			from := uids[i2]
			to := uids[rand.Int63n(int64(userCount))]
			m1.Lock()
			msgTo[from] = to
			m1.Unlock()

			chat, err := dao.ChatDao.CreateChat(dao.ChatTypeUser, from, to)
			if err != nil {
				panic(err)
			}
			m2.Lock()
			cids[from] = chat.Cid
			m2.Unlock()

			userChat, err := dao.ChatDao.NewUserChat(chat.Cid, from, to, dao.ChatTypeUser)
			if err != nil {
				panic(err)
			}
			m3.Lock()
			ucIds[from] = userChat.UcId
			m3.Unlock()
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
		cids[from] = 1
		ucIds[from] = 1
	}
}

func sleepRndMilleSec(start int32, end int32) {
	n := rand.Int31n(end - start)
	n = start + n
	time.Sleep(time.Duration(n) * time.Millisecond)
}
