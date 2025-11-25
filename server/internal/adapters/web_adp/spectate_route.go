package web_adp

import (
	"github.com/gin-gonic/gin"
	"github.com/vanshjangir/vango/server/internal/applications/util_app"
)

func (gh *GinHandler) spectate(ctx *gin.Context) {
	wsUrl := util_app.Pick();
	ctx.JSON(200, gin.H{"wsurl": wsUrl});
}
