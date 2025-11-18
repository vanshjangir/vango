package ws_adp

import "github.com/gorilla/websocket"

type WsGameRepo struct {
	ws *websocket.Conn
}

func (r *WsGameRepo) Send(data []byte) error {
	return r.ws.WriteMessage(websocket.TextMessage, data)
}

func (r *WsGameRepo) Receive() ([]byte, error) {
	_, data, err := r.ws.ReadMessage()
	return data, err
}

func (r *WsGameRepo) Close() error {
	return r.ws.Close()
}

func NewWebsocketGameRepo(ws *websocket.Conn) *WsGameRepo {
	return &WsGameRepo{ws: ws}
}
