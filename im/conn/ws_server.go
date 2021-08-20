package conn

import (
	"fmt"
	"go_im/im/comm"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WsServerOptions struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type WsServer struct {
	options  *WsServerOptions
	upgrader websocket.Upgrader
	handler  ConnectionHandler
}

// NewWsServer options can be nil, use default value when nil.
func NewWsServer(options *WsServerOptions) *WsServer {

	if options == nil {
		options = &WsServerOptions{
			Host:         "0.0.0.0",
			Port:         8080,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
	}
	ws := new(WsServer)
	ws.options = options
	ws.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 65536,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return ws
}

func (ws *WsServer) handleWebSocketRequest(writer http.ResponseWriter, request *http.Request) {

	conn, err := ws.upgrader.Upgrade(writer, request, nil)
	if err != nil {
		comm.Slog.E("upgrade http to ws error", err)
		return
	}

	proxy := ConnectionProxy{
		conn: NewWsConnection(conn, ws.options),
	}
	ws.handler(proxy)
}

func (ws *WsServer) Handler(handler ConnectionHandler) {
	ws.handler = handler
}

func (ws *WsServer) Run() {

	http.HandleFunc("/ws", ws.handleWebSocketRequest)

	addr := fmt.Sprintf("%s:%d", ws.options.Host, ws.options.Port)
	fmt.Printf("websocket run on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}

}
