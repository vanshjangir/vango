package main

import (
	"fmt"

	"github.com/vanshjangir/rapidgo/server/internal/adapters/postgres_adp"
	"github.com/vanshjangir/rapidgo/server/internal/adapters/web_adp"
	"github.com/vanshjangir/rapidgo/server/internal/applications/mm_app"
	"github.com/vanshjangir/rapidgo/server/internal/applications/user_app"
)

func main() {
	db := postgres_adp.SetupDB()
	userRepo := postgres_adp.NewPostgresUserRepo(db)
    userService := user_app.NewUserService(userRepo)
	matchMakingService := mm_app.NewMatchMakingService(userRepo, nil)

	httpHandler := web_adp.NewGinHandler(userService, nil, matchMakingService)
	httpHandler.RegisterRoutes()
	httpHandler.Run()

	fmt.Println("Starting server")
}
