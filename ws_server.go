package poker

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WSServer struct {
	*websocket.Conn
}

func NewWSServer(w http.ResponseWriter, r *http.Request) *WSServer {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("problem upgrading connection to WebSockets %v\n", err)
	}

	return &WSServer{conn}
}

func (w *WSServer) WaitForMsg() string {
	_, msg, err := w.ReadMessage()
	if err != nil {
		log.Printf("error reading frmo websocket %v", err)
	}

	return string(msg)
}

func (w *WSServer) Write(p []byte) (n int, err error) {
	err = w.WriteMessage(websocket.TextMessage, p)

	if err != nil {
		return 0, err
	}

	return len(p), nil
}
