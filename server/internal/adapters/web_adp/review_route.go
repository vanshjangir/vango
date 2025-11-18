package web_adp

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (gh *GinHandler) review(ctx *gin.Context) {
	gameid, err := strconv.Atoi(ctx.Query("gameid"))
	if err != nil {
		ctx.Status(400)
		ctx.JSON(400, gin.H{"error": "Wrong Game Id"})
		return
	}

	review, err := gh.gs.Review(gameid)
	if err != nil {
		ctx.Status(400)
		ctx.JSON(400, gin.H{"error": "Review not found"})
	} else {
		ctx.JSON(200, review)
	}
}
