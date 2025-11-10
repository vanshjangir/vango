package ws_adp

import (
	"log"
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

	//username, err := wsGameService.Auth(wsh.us)
	//if err != nil {
	//	wsGameService.Send([]byte("Authentication unsuccessfull"))
	//	wsGameService.Close()
	//	return
	//}
	username := ctx.Query("username")

	game, err := wsGameService.SetupGame(username)
	if err != nil {
		log.Println("Error setting up game:", err)
		wsGameService.Send([]byte("Server Error occurred"))
		return
	}

	err = wsGameService.SendStartConfirmation(game);
	if err != nil {
		log.Println("Error sending start confirmation:", err)
		wsGameService.Close()
		return
	}

	go wsGameService.Play(game)
}
