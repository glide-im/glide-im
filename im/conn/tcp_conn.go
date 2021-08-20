package conn

type TcpConnection struct {
}

func (t *TcpConnection) Write(message Serializable) error {
	panic("implement me")
}

func (t *TcpConnection) Read(message Serializable) error {
	panic("implement me")
}

func (t *TcpConnection) Close() error {
	panic("implement me")
}
