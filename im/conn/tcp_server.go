package conn

import "net"

type TcpServer struct {
	handler ConnectionHandler
}

func NewTcpServer() *TcpServer {
	return &TcpServer{}
}

func (t *TcpServer) SetConnHandler(handler ConnectionHandler) {
	t.handler = handler
}

func (t *TcpServer) Run(host string, port int) error {
	tcp, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.ParseIP(host),
		Port: port,
	})
	if err != nil {
		return err
	}
	for {
		acceptTCP, err := tcp.AcceptTCP()
		if err != nil {
			return err
		}
		conn := ConnectionProxy{
			conn: NewTcpConn(acceptTCP),
		}
		t.handler(conn)
	}
}
