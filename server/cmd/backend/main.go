package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/vanshjangir/vango/server/internal/adapters/postgres_adp"
	"github.com/vanshjangir/vango/server/internal/adapters/pubsub_adp"
	"github.com/vanshjangir/vango/server/internal/adapters/web_adp"
	"github.com/vanshjangir/vango/server/internal/applications/game_app"
	"github.com/vanshjangir/vango/server/internal/applications/mm_app"
	"github.com/vanshjangir/vango/server/internal/applications/user_app"
	"github.com/vanshjangir/vango/server/internal/applications/util_app"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := util_app.Load("./gameservers.json"); err != nil {
		log.Fatal("Error loading game servers", err)
	}

	db := postgres_adp.SetupDB()

	userRepo := postgres_adp.NewPostgresUserRepo(db)
	gameRepo := postgres_adp.NewPostgresGameRepo(db)
	pubsubRepo := pubsub_adp.NewPubsubRepo()

	userService := user_app.NewUserService(userRepo)
	gameService := game_app.NewGameService(gameRepo)
	matchMakingService := mm_app.NewMatchMakingService(userRepo, gameRepo, pubsubRepo)

	httpHandler := web_adp.NewGinHandler(userService, gameService, matchMakingService)
	httpHandler.RegisterRoutes()
	httpHandler.Run()

	log.Println("Starting server")
}
