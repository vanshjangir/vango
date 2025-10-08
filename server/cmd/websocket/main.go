package main

import (
	"fmt"

	"github.com/vanshjangir/rapidgo/server/internal/adapters/postgres_adp"
	"github.com/vanshjangir/rapidgo/server/internal/adapters/ws_adp"
	"github.com/vanshjangir/rapidgo/server/internal/applications/game_app"
	"github.com/vanshjangir/rapidgo/server/internal/applications/user_app"
)

func main() {
	db := postgres_adp.SetupDB()
	
    userRepo := postgres_adp.NewPostgresUserRepo(db)
    gameRepo := postgres_adp.NewPostgresGameRepo(db)
    
	userService := user_app.NewUserService(userRepo)
    gameService := game_app.NewGameService(gameRepo)

	wsHandler := ws_adp.NewWsHandler(userService, gameService)
	wsHandler.RegisterRoutes()
	wsHandler.Run()

	fmt.Println("Starting server")
}
