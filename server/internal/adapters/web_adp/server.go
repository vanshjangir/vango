package web_adp

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vanshjangir/rapidgo/server/internal/ports"
)

type GinHandler struct {
	r  *gin.Engine
	us ports.UserService
	gs ports.GameService
	ms ports.MatchMakingService
}

func NewGinHandler(
	us ports.UserService,
	gs ports.GameService,
	ms ports.MatchMakingService,
) *GinHandler {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	return &GinHandler{r: r, us: us, gs: gs, ms: ms}
}

func (gh *GinHandler) RegisterRoutes() {
	gh.r.GET("/profile", gh.profile)
	gh.r.GET("/review", gh.review)
	gh.r.POST("/login", gh.login)
	gh.r.POST("/signup", gh.signup)

	gh.r.GET("/findgame", gh.findGame)
	gh.r.POST("/changeusername", gh.httpAuth, gh.changeUsername)
}

func (gh *GinHandler) Run() {
	gh.r.Run()
}
