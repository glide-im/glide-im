package conn

import "net"

type TcpConnection struct {
	c *net.TCPConn
}

func NewTcpConn(c *net.TCPConn) *TcpConnection {
	return &TcpConnection{c: c}
}

func (t TcpConnection) Write(data []byte) error {
	_, err := t.c.Write(data)
	return err
}

func (t TcpConnection) Read() ([]byte, error) {
	var b []byte
	_, err := t.c.Read(b)
	return b, err
}

func (t TcpConnection) Close() error {
	return t.c.Close()
}

func (t TcpConnection) GetConnInfo() *ConnectionInfo {
	addr := t.c.RemoteAddr().(*net.TCPAddr)
	return &ConnectionInfo{
		Ip:   addr.IP.String(),
		Port: addr.Port,
		Addr: t.c.RemoteAddr().String(),
	}
}
