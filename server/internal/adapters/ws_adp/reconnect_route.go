package ws_adp

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func (wsh *WsHandler) reconnect(ctx *gin.Context) {
	w, r := ctx.Writer, ctx.Request
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Could not able to upgrade connection"})
		return
	}

	wsGameRepo := NewWebsocketGameRepo(c)

	//username, err := wsGameService.Auth(wsh.us)
	//if err != nil {
	//	wsGameService.Send([]byte("Authentication unsuccessfull"))
	//	wsGameService.Close()
	//	return
	//}
	username := ctx.Query("username")
	game, err := wsh.ws.LoadExistingGame(username, wsGameRepo)
	if err != nil {
		log.Println("Error reconnecting to game for user", username, err)
		c.WriteMessage(websocket.TextMessage, []byte("Server Error occurred"))
		c.Close()
		return
	}
	game.IsOnline = true
	*game.ReconnectChan <- true
	log.Printf("Player %v reconnected\n", username)
}
