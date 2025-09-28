package webAdp

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (gh *GinHandler) httpAuth(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
    authData := strings.Split(authHeader, " ")
    tokenType := authData[0]
    token := authData[1]

    if tokenType == "" || token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token not found"})
        ctx.Abort()
        return
    }

    var err error
    var username string
    switch tokenType {
    case "google":
        username, err = gh.us.AuthGoogle(token)
    case "guest":
        username, err = gh.us.AuthGuest(token)
    default:
        err = fmt.Errorf("Unsupported token type")
    }

    if err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
        ctx.Abort()
    } else {
        ctx.Set("username", username)
        ctx.Next()
    }
}
