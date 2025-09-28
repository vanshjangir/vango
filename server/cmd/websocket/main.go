package main

import (
	"fmt"

	"github.com/vanshjangir/rapidgo/server/internal/adapters/postgresAdp"
	"github.com/vanshjangir/rapidgo/server/internal/adapters/wsAdp"
	"github.com/vanshjangir/rapidgo/server/internal/applications"
)

func main() {
	db := postgresAdp.SetupDB()
	
    userRepo := postgresAdp.NewPostgresUserRepo(db)
    userService := applications.NewUserService(userRepo)

	wsHandler := wsAdp.NewWsHandler(userService, nil)
	wsHandler.RegisterRoutes()
	wsHandler.Run()

	fmt.Println("Starting server")
}
