package wsAdp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/vanshjangir/rapidgo/server/internal/applications"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (wh *WsHandler) play(ctx *gin.Context) {
	w, r := ctx.Writer, ctx.Request
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Could not able to upgrade connection"})
        return
	}
	
	wsGameRepo := NewWebsocketGameRepo(c)
    wsGameService := applications.NewWsGameService(wsGameRepo)
	wsGameService.Receive()
	username, err := wsGameService.Auth(wh.us)
	if err != nil {
		wsGameService.Send([]byte("Authentication unsuccessfull"))
		wsGameService.Close()
		return
	}

	wsGameService.Play(username)
}
