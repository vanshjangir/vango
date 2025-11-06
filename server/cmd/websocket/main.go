package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/vanshjangir/rapidgo/server/internal/adapters/postgres_adp"
	"github.com/vanshjangir/rapidgo/server/internal/adapters/pubsub_adp"
	"github.com/vanshjangir/rapidgo/server/internal/adapters/ws_adp"
	"github.com/vanshjangir/rapidgo/server/internal/applications/game_app"
	"github.com/vanshjangir/rapidgo/server/internal/applications/user_app"
	"github.com/vanshjangir/rapidgo/server/internal/applications/ws_app"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	
	db := postgres_adp.SetupDB()
	
    userRepo := postgres_adp.NewPostgresUserRepo(db)
    gameRepo := postgres_adp.NewPostgresGameRepo(db)
    wsGameRepo := ws_adp.NewWebsocketGameRepo(nil)
    pubsubRepo := pubsub_adp.NewPubsubRepo()
    
	userService := user_app.NewUserService(userRepo)
    gameService := game_app.NewGameService(gameRepo)
    wsGameService := ws_app.NewWsGameService(wsGameRepo, pubsubRepo, gameRepo, userRepo)

	wsHandler := ws_adp.NewWsHandler(userService, gameService, wsGameService)
	wsHandler.RegisterRoutes()
	wsHandler.Run()

	log.Println("Starting server")
}
