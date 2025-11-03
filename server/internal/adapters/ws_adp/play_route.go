package ws_adp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (wsh *WsHandler) play(ctx *gin.Context) {
	w, r := ctx.Writer, ctx.Request
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Could not able to upgrade connection"})
        return
	}

	wsGameRepo := NewWebsocketGameRepo(c)
    wsGameService := wsh.ws.CopyWsGameService(wsGameRepo)

	username, err := wsGameService.Auth(wsh.us)
	if err != nil {
		wsGameService.Send([]byte("Authentication unsuccessfull"))
		wsGameService.Close()
		return
	}

	game, err := wsGameService.SetupGame(username)
	if err != nil {
		wsGameService.Send([]byte("Server Error occurred"))
		return
	}

	err = wsGameService.SendStartConfirmation();
	if err != nil {
		wsGameService.Close()
		return
	}

	go wsGameService.Play(game)
}
