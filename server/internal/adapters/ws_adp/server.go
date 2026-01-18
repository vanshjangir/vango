package ws_adp

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vanshjangir/vango/server/internal/ports"
)

type WsHandler struct {
	r  *gin.Engine
	us ports.UserService
	gs ports.GameService
	ws ports.WsGameService
	ss ports.SpectateService
}

func NewWsHandler(
	us ports.UserService,
	gs ports.GameService,
	ws ports.WsGameService,
	ss ports.SpectateService,
) *WsHandler {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	return &WsHandler{r: r, us: us, gs: gs, ws: ws, ss: ss}
}

func (wsh *WsHandler) RegisterRoutes() {
	wsh.r.GET("/play", wsh.wsAuth, wsh.play)
	wsh.r.GET("/spectate", wsh.wsAuth, wsh.spectate)
}

func (wsh *WsHandler) Run() {
	wsh.r.Run()
}
