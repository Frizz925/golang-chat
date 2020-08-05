package handlers

import (
	"golang-chat/lib"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	stream *lib.Stream
}

var _ http.Handler = (*WebSocketHandler)(nil)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewWebSocketHandler(stream *lib.Stream) *WebSocketHandler {
	return &WebSocketHandler{stream}
}

func (h *WebSocketHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Print(err)
	} else {
		h.ServeWebSocket(conn)
	}
}

func (h *WebSocketHandler) ServeWebSocket(conn *websocket.Conn) {
	defer conn.Close()
	ch, id := h.stream.Subscribe()
	defer h.stream.Unsubscribe(id)
	for {
		message := <-ch
		if message == "" {
			break
		}
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			if err != io.EOF {
				log.Print(err)
			}
			break
		}
	}
}
