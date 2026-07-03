package gate

import (
	"errors"
	"github.com/glide-im/glide/pkg/conn"
	"github.com/glide-im/glide/pkg/logger"
	"github.com/glide-im/glide/pkg/messages"
	"github.com/glide-im/glide/pkg/timingwheel"
	"sync"
	"sync/atomic"
	"time"
)

// tw is a timer for heartbeat.
var tw = timingwheel.NewTimingWheel(time.Millisecond*500, 3, 20)

const (
	defaultServerHeartbeatDuration = time.Second * 30
	defaultHeartbeatDuration       = time.Second * 20
	defaultHeartbeatLostLimit      = 3
	defaultCloseImmediately        = false
)

// client state
const (
	_ int32 = iota
	// stateRunning client is running, can runRead and runWrite message.
	stateRunning
	// stateClosed client is closed, cannot do anything.
	stateClosed
)

// ClientConfig client config
type ClientConfig struct {

	// ClientHeartbeatDuration is the duration of heartbeat.
	ClientHeartbeatDuration time.Duration

	// ServerHeartbeatDuration is the duration of server heartbeat.
	ServerHeartbeatDuration time.Duration

	// HeartbeatLostLimit is the max lost heartbeat count.
	HeartbeatLostLimit int

	// CloseImmediately true express when client exit, discard all message in queue, and close connection immediately,
	// otherwise client will close runRead, and mark as stateClosing, the client cannot receive and enqueue message,
	// after all message in queue is sent, client will close runWrite and connection.
	CloseImmediately bool
}

type MessageInterceptor = func(dc DefaultClient, msg *messages.GlideMessage) bool

type DefaultClient interface {
	Client

	SetCredentials(credentials *ClientAuthCredentials)

	GetCredentials() *ClientAuthCredentials

	AddMessageInterceptor(interceptor MessageInterceptor)
}

var _ DefaultClient = (*UserClient)(nil)

// UserClient represent a user conn client.
type UserClient struct {

	// conn is the real connection
	conn conn.Connection

	// state is the client state
	state int32

	// queuedMessage message count in the messages channel
	queuedMessage int64
	// messages is the buffered channel for message to push to client.
	messages chan *messages.GlideMessage

	// closeReadCh is the channel for runRead goroutine to close
	closeReadCh chan struct{}
	// closeWriteCh is the channel for runWrite goroutine to close
	closeWriteCh chan struct{}

	// closeWriteOnce is the once for close runWrite goroutine
	closeWriteOnce sync.Once
	// closeReadOnce is the once for close runRead goroutine
	closeReadOnce sync.Once
	// closeOnce is the once for close messages channel and connection
	closeOnce sync.Once

	// hbC is the timer for client heartbeat
	hbC *timingwheel.Task
	// hbS is the timer for server heartbeat
	hbS *timingwheel.Task
	// hbLost is the count of heartbeat lost
	hbLost int32

	// info is the client info
	info *Info

	mu          sync.Mutex
	credentials *ClientAuthCredentials

	// mgr the client manager which manage this client
	mgr Gateway
	// msgHandler client message handler
	msgHandler MessageHandler

	// config is the client config
	config *ClientConfig
}

func NewClientWithConfig(conn conn.Connection, mgr Gateway, handler MessageHandler, config *ClientConfig) DefaultClient {
	if config == nil {
		config = &ClientConfig{
			ClientHeartbeatDuration: defaultHeartbeatDuration,
			ServerHeartbeatDuration: defaultServerHeartbeatDuration,
			HeartbeatLostLimit:      defaultHeartbeatLostLimit,
			CloseImmediately:        false,
		}
	}

	ret := UserClient{
		conn:         conn,
		messages:     make(chan *messages.GlideMessage, 100),
		closeReadCh:  make(chan struct{}),
		closeWriteCh: make(chan struct{}),
		hbC:          tw.After(config.ClientHeartbeatDuration),
		hbS:          tw.After(config.ServerHeartbeatDuration),
		info: &Info{
			ConnectionAt: time.Now().UnixMilli(),
			CliAddr:      conn.GetConnInfo().Addr,
		},
		mgr:        mgr,
		msgHandler: handler,
		config:     config,
	}
	return &ret
}

func NewClient(conn conn.Connection, mgr Gateway, handler MessageHandler) DefaultClient {
	return NewClientWithConfig(conn, mgr, handler, nil)
}

func (c *UserClient) SetCredentials(credentials *ClientAuthCredentials) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.credentials = credentials
	c.info.ConnectionId = credentials.ConnectionID
	if credentials.ConnectionConfig != nil {
		c.config.HeartbeatLostLimit = credentials.ConnectionConfig.AllowMaxHeartbeatLost
		c.config.CloseImmediately = credentials.ConnectionConfig.CloseImmediately
		c.config.ClientHeartbeatDuration = time.Duration(credentials.ConnectionConfig.HeartbeatDuration) * time.Second
	}
}

func (c *UserClient) GetCredentials() *ClientAuthCredentials {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.credentials
}

func (c *UserClient) AddMessageInterceptor(interceptor MessageInterceptor) {
	h := c.msgHandler
	c.msgHandler = func(cliInfo *Info, msg *messages.GlideMessage) {
		if interceptor(c, msg) {
			return
		}
		h(cliInfo, msg)
	}
}

func (c *UserClient) GetInfo() Info {
	return *c.info
}

// SetID set client id.
func (c *UserClient) SetID(id ID) {
	c.info.ID = id
}

// IsRunning return true if client is running
func (c *UserClient) IsRunning() bool {
	return atomic.LoadInt32(&c.state) == stateRunning
}

// EnqueueMessage enqueue message to client message queue.
func (c *UserClient) EnqueueMessage(msg *messages.GlideMessage) (err error) {
	if atomic.LoadInt32(&c.state) == stateClosed {
		return errors.New("client has closed")
	}
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("client has closed")
		}
	}()
	logger.I("EnqueueMessage ID=%s msg=%v", c.info.ID, msg)
	select {
	case c.messages <- msg:
		atomic.AddInt64(&c.queuedMessage, 1)
	default:
		logger.E("msg chan is full, id=%v", c.info.ID)
	}
	return nil
}

// runRead message from client.
func (c *UserClient) runRead() {
	defer func() {
		err := recover()
		if err != nil {
			logger.E("read message panic: %v", err)
			c.Exit()
		}
	}()

	readChan, done := messageReader.ReadCh(c.conn)
	var closeReason string
	for {
		select {
		case <-c.closeReadCh:
			if closeReason == "" {
				closeReason = "closed initiative"
			}
			goto STOP
		case <-c.hbC.C:
			if !c.IsRunning() {
				goto STOP
			}
			atomic.AddInt32(&c.hbLost, 1)
			if int(atomic.LoadInt32(&c.hbLost)) > c.config.HeartbeatLostLimit {
				closeReason = "heartbeat lost"
				c.Exit()
			}
			c.hbC.Cancel()
			c.hbC = tw.After(c.config.ClientHeartbeatDuration)
			_ = c.EnqueueMessage(messages.NewMessage(0, messages.ActionHeartbeat, nil))
		case msg := <-readChan:
			if msg == nil {
				closeReason = "readCh closed"
				c.Exit()
				continue
			}
			if msg.err != nil {
				if messages.IsDecodeError(msg.err) {
					_ = c.EnqueueMessage(messages.NewMessage(0, messages.ActionNotifyError, msg.err.Error()))
					continue
				}
				closeReason = msg.err.Error()
				c.Exit()
				continue
			}
			if c.info.ID == "" {
				closeReason = "client not logged"
				c.Exit()
				break
			}
			atomic.StoreInt32(&c.hbLost, 0)
			c.hbC.Cancel()
			c.hbC = tw.After(c.config.ClientHeartbeatDuration)

			if msg.m.GetAction() == messages.ActionHello {
				c.handleHello(msg.m)
			} else {
				c.msgHandler(c.info, msg.m)
			}
			msg.Recycle()
		}
	}
STOP:
	close(done)
	c.hbC.Cancel()
	logger.I("read exit, reason=%s", closeReason)
}

// runWrite message to client.
func (c *UserClient) runWrite() {
	defer func() {
		err := recover()
		if err != nil {
			logger.D("write message error, exit client: %v", err)
			c.Exit()
		}
	}()

	var closeReason string
	for {
		select {
		case <-c.closeWriteCh:
			if closeReason == "" {
				closeReason = "closed initiative"
			}
			goto STOP
		case <-c.hbS.C:
			if !c.IsRunning() {
				closeReason = "client not running"
				goto STOP
			}
			_ = c.EnqueueMessage(messages.NewMessage(0, messages.ActionHeartbeat, nil))
			c.hbS.Cancel()
			c.hbS = tw.After(c.config.ServerHeartbeatDuration)
		case m := <-c.messages:
			if m == nil {
				closeReason = "message is nil, maybe client has closed"
				c.Exit()
				break
			}
			c.write2Conn(m)
			c.hbS.Cancel()
			c.hbS = tw.After(c.config.ServerHeartbeatDuration)
		}
	}
STOP:
	c.hbS.Cancel()
	logger.D("write exit, addr=%s, reason:%s", c.info.CliAddr, closeReason)
}

// Exit client, note: exit client will not close conn right now, but will close when message chan is empty.
// It's close read right now, and close write2Conn when all message in queue is sent.
func (c *UserClient) Exit() {
	if !atomic.CompareAndSwapInt32(&c.state, stateRunning, stateClosed) {
		return
	}

	id := c.info.ID
	// exit by client self, remove client from manager
	if c.mgr != nil && id != "" {
		_ = c.mgr.ExitClient(id)
	}
	c.SetID("")
	c.mgr = nil
	c.stopReadWrite()

	if c.config.CloseImmediately {
		// dropping all message in queue and close connection immediately
		c.close()
	} else {
		// close connection when all message in queue is sent
		go func() {
			for {
				select {
				case m := <-c.messages:
					c.write2Conn(m)
				default:
					goto END
				}
			}
		END:
			c.close()
		}()
	}
}

func (c *UserClient) Run() {
	logger.I("new client running addr:%s id:%s", c.conn.GetConnInfo().Addr, c.info.ID)

	c.messages = make(chan *messages.GlideMessage, 100)
	c.closeReadCh = make(chan struct{})
	c.closeWriteCh = make(chan struct{})
	c.closeWriteOnce = sync.Once{}
	c.closeReadOnce = sync.Once{}
	c.closeOnce = sync.Once{}
	atomic.StoreInt32(&c.hbLost, 0)

	atomic.StoreInt32(&c.state, stateRunning)

	go c.runRead()
	go c.runWrite()
}

func (c *UserClient) isClosed() bool {
	return atomic.LoadInt32(&c.state) == stateClosed
}

func (c *UserClient) close() {
	c.closeOnce.Do(func() {
		close(c.messages)
		_ = c.conn.Close()
	})
}

func (c *UserClient) write2Conn(m *messages.GlideMessage) {
	b, err := codec.Encode(m)
	if err != nil {
		logger.E("serialize output message", err)
		return
	}
	err = c.conn.Write(b)
	atomic.AddInt64(&c.queuedMessage, -1)
	if err != nil {
		logger.D("runWrite error: %s", err.Error())
		c.closeWriteOnce.Do(func() {
			close(c.closeWriteCh)
		})
	}
}

func (c *UserClient) stopReadWrite() {
	c.closeWriteOnce.Do(func() {
		close(c.closeWriteCh)
	})
	c.closeReadOnce.Do(func() {
		close(c.closeReadCh)
	})
}

func (c *UserClient) handleHello(m *messages.GlideMessage) {
	hello := messages.Hello{}
	err := m.Data.Deserialize(&hello)
	if err != nil {
		_ = c.EnqueueMessage(messages.NewMessage(0, messages.ActionNotifyError, "invalid handleHello message"))
	} else {
		c.info.Version = hello.ClientVersion
	}
}
