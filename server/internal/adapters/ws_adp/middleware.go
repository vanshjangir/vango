package ws_adp

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (wsh *WsHandler) wsAuth(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Sec-Websocket-Protocol")
	authData := strings.Split(authHeader, ".")
	tokenType := authData[0]
	token := strings.Join(authData[1:], ".")

	if tokenType == "" || token == "" {
		log.Println("Token not found")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token not found"})
		ctx.Abort()
		return
	}

	var err error
	var username string
	switch tokenType {
	case "google":
		username, err = wsh.us.AuthGoogle(token)
	case "guest":
		username, err = wsh.us.AuthGuest(token)
	default:
		err = fmt.Errorf("Unsupported token type")
	}

	if err != nil {
		log.Println("Authentication failed:", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		ctx.Abort()
	} else {
		ctx.Set("username", username)
		ctx.Set("usertype", tokenType)
		ctx.Set("protocolHeader", authHeader)
		ctx.Next()
	}
}
