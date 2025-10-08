package web_adp

import "github.com/gin-gonic/gin"

func (gh *GinHandler) review(ctx *gin.Context) {
	gameid := ctx.Query("gameid")
    gh.gs.Review(gameid)
}
