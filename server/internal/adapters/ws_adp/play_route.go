package ws_adp

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func getUpgrader (ctx *gin.Context) (websocket.Upgrader, error) {
	var upgrader websocket.Upgrader
	protocolHeaderItf, exists := ctx.Get("protocolHeader")
	if !exists {
		return upgrader, fmt.Errorf("getUpgrader: protocolHeader not found")
	}

	protocolHeader := protocolHeaderItf.(string)
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{protocolHeader},
	}

	return upgrader, nil
}

func (wsh *WsHandler) startNewGame(
	username string,
	wsGameRepo *WsGameRepo,
	c *websocket.Conn,
) {
	game, err := wsh.ws.SetupGame(username, wsGameRepo)
	if err != nil {
		log.Println("Error setting up game:", err)
		c.WriteMessage(websocket.TextMessage, []byte("Server Error occurred"))
		c.Close()
		return
	}
	
	go wsh.ws.Play(game)
}

func (wsh *WsHandler) reconnectExistingGame(
	username string,
	wsGameRepo *WsGameRepo,
	c *websocket.Conn,
) {
	game, err := wsh.ws.LoadExistingGame(username, wsGameRepo)
	if err != nil {
		log.Println("Error reconnecting to game for user", username, err)
		c.WriteMessage(websocket.TextMessage, []byte("Server Error occurred"))
		c.Close()
		return
	}
	
	game.IsOnline = true
	*game.ReconnectChan <- true
	
	err = wsh.ws.SendStartConfirmation(game)
	if err != nil {
		log.Println("reconnectExistingGame: SendStartConfirmation:", username, err)
		c.WriteMessage(websocket.TextMessage, []byte("Server Error occurred"))
		c.Close()
		return
	}
	
	log.Printf("Player %v reconnected\n", username)
}

func (wsh *WsHandler) play(ctx *gin.Context) {
	usernameItf, exists := ctx.Get("username")
	if !exists {
		log.Println("username and usertype not found")
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	
	upgrader, err := getUpgrader(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Could not able to get upgrader"})
		return
	}
	
	username := usernameItf.(string)
	w, r := ctx.Writer, ctx.Request
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection", err)
		ctx.JSON(500, gin.H{"error": "Could not able to upgrade connection"})
		c.Close()
		return
	}
	
	wsGameRepo := NewWebsocketGameRepo(c)
	
	if wsh.ws.GameExists(username) {
		wsh.reconnectExistingGame(username, wsGameRepo, c)
	} else {
		wsh.startNewGame(username, wsGameRepo, c)
	}
}
