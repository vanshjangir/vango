package ws_adp

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vanshjangir/rapidgo/server/internal/ports"
)

type WsHandler struct {
	r  *gin.Engine
	us ports.UserService
	gs ports.GameService
	ws ports.WsGameService
}

func NewWsHandler(
	us ports.UserService,
	gs ports.GameService,
    ws ports.WsGameService,
) *WsHandler {
    r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	return &WsHandler{r: r, us: us, gs: gs, ws: ws}
}

func (wsh *WsHandler) RegisterRoutes() {
	wsh.r.GET("/play", wsh.play)
}

func (wsh *WsHandler) Run() {
	wsh.r.Run()
}
