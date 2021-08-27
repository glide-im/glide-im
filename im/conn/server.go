package conn

type ConnectionHandler func(conn Connection)

type Server interface {
	SetConnHandler(handler ConnectionHandler)
	Run(host string, port int) error
}
