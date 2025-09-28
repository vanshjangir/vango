package wsAdp

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vanshjangir/rapidgo/server/internal/ports"
)

type WsHandler struct {
	r  *gin.Engine
	us ports.UserService
	gs ports.GameService
}

func NewWsHandler(
    us ports.UserService,
    gs ports.GameService,
) *WsHandler {
    r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	return &WsHandler{r: r, us: us, gs: gs}
}

func (wh *WsHandler) RegisterRoutes() {
	wh.r.GET("/play", wh.play)
}

func (wh *WsHandler) Run() {
	wh.r.Run()
}
