package client

import (
	"go_im/im/conn"
	"go_im/im/message"
	"sync"
)

var messageReader MessageReader

//var codec message.Codec = message.ProtobufCodec{}
var codec message.Codec = message.DefaultCodec

// recyclePool 回收池, 减少临时对象, 回收复用 readerRes
var recyclePool sync.Pool

func init() {
	recyclePool = sync.Pool{
		New: func() interface{} {
			return &readerRes{}
		},
	}
	SetMessageReader(&defaultReader{})
}

func SetMessageReader(s MessageReader) {
	messageReader = s
}

type readerRes struct {
	err error
	m   *message.Message
}

// Recycle 回收当前对象, 一定要在用完后调用这个方法, 否则无法回收
func (r *readerRes) Recycle() {
	r.m = nil
	r.err = nil
	recyclePool.Put(r)
}

// MessageReader 表示一个从连接中(Connection)读取消息的读取者, 可以用于定义如何从连接中读取并解析消息.
type MessageReader interface {

	// Read 阻塞读取, 会阻塞当前协程
	Read(conn conn.Connection) (*message.Message, error)

	// ReadCh 返回两个管道, 第一个用于读取内容, 第二个用于发送停止读取, 停止读取时切记要发送停止信号
	ReadCh(conn conn.Connection) (<-chan *readerRes, chan<- struct{})
}

type defaultReader struct{}

func (d *defaultReader) ReadCh(conn conn.Connection) (<-chan *readerRes, chan<- struct{}) {
	c := make(chan *readerRes, 5)
	done := make(chan struct{})
	// TODO 2021-12-3 IO重构
	go func() {
		for {
			select {
			case <-done:
				goto CLOSE
			default:
				m, err := d.Read(conn)
				res := recyclePool.Get().(*readerRes)
				res.err = err
				if err != nil {
					c <- res
					goto DONE
				}
				res.m = m
				c <- res
			}
		}
	CLOSE:
		close(done)
	DONE:
		close(c)
	}()
	return c, done
}

func (d *defaultReader) Read(conn conn.Connection) (*message.Message, error) {
	// TODO 2021-12-3 校验数据包
	bytes, err := conn.Read()
	if err != nil {
		return nil, err
	}
	m := message.NewEmptyMessage()
	err = codec.Decode(bytes, m)
	return m, err
}
