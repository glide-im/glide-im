package conn

type TcpConnection struct {
}

func (t TcpConnection) Write(data []byte) error {
	panic("implement me")
}

func (t TcpConnection) Read() ([]byte, error) {
	panic("implement me")
}

func (t TcpConnection) Close() error {
	panic("implement me")
}
