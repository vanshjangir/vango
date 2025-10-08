package web_adp

import (
	"github.com/gin-gonic/gin"
	"github.com/vanshjangir/rapidgo/server/internal/domain"
)

func (gh *GinHandler) findGame(ctx *gin.Context) {
	var wp domain.WaitingPlayer
	if err := ctx.ShouldBindJSON(&wp); err != nil {
		ctx.JSON(400, gin.H{"error": "invalid payload"})
		return
	}

    if wsurl, err := gh.ms.Match(wp); err != nil {
        ctx.JSON(500, gin.H{"error": "Matching failed"})
    } else {
        ctx.JSON(200, gin.H{"wsurl": wsurl})
    }
}
