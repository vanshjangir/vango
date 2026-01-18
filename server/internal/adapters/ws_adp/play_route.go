package ws_adp

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/vanshjangir/vango/server/internal/domain"
)

func getUpgrader(ctx *gin.Context) (websocket.Upgrader, error) {
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

func (wsh *WsHandler) spectateGame(gameId int, c *websocket.Conn) {
	game := wsh.ws.GetGameFromId(gameId)
	if game == nil {
		log.Println("Game not found", gameId)
		c.WriteMessage(websocket.TextMessage, []byte("Game not found"))
		c.Close()
		return
	}

	wsGameRepo := NewWebsocketGameRepo(c)
	wsh.ss.AddSpectator(game, wsGameRepo)

	var blackGame, whiteGame *domain.Game
	if game.Color == domain.BlackColor {
		blackGame = game
		whiteGame = wsh.ws.GetGameFromPlayerName(game.OpName)
	} else {
		whiteGame = game
		blackGame = wsh.ws.GetGameFromPlayerName(game.OpName)
	}

	err := wsh.ss.SendStartConfirmation(wsGameRepo)
	if err != nil {
		log.Println("spectateGame: SendSyncState:", err)
		return
	}
	
	err = wsh.ss.SendSyncState(blackGame, whiteGame, wsGameRepo)
	if err != nil {
		log.Println("spectateGame: SendSyncState:", err)
		return
	}
	log.Println("New spectator added to game:", gameId)
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

	gameIdString := ctx.Param("gameid")
	if gameIdString != "" {
		gameId, err := strconv.Atoi(ctx.Param("gameid"))
		if wsh.ss.IsPlayerInGame(username, gameId) == false {
			if err != nil {
				log.Println("Invalid gameid", ctx.Param("gameid"))
			}
			wsh.spectateGame(gameId, c)
			return
		}
	}

	if wsh.ws.GameExists(username) {
		wsh.reconnectExistingGame(username, wsGameRepo, c)
	} else {
		wsh.startNewGame(username, wsGameRepo, c)
	}
}
