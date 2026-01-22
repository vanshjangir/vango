package ws_adp

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/vanshjangir/vango/server/internal/domain"
)

func (wsh *WsHandler) spectate(ctx *gin.Context) {
	upgrader, err := getUpgrader(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Could not able to get upgrader"})
		return
	}

	w, r := ctx.Writer, ctx.Request
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection", err)
		ctx.JSON(500, gin.H{"error": "Could not able to upgrade connection"})
		c.Close()
		return
	}

	gameId, err := strconv.Atoi(ctx.Query("gameid"))
	if err != nil {
		log.Println("Invalid gameid", ctx.Query("gameid"))
		c.WriteMessage(websocket.TextMessage, []byte("Invalid game id"))
		c.Close()
		return
	}

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

	err = wsh.ss.SendStartConfirmation(wsGameRepo)
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
