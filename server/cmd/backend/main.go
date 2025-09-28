package main

import (
	"fmt"

	"github.com/vanshjangir/rapidgo/server/internal/adapters/postgresAdp"
	"github.com/vanshjangir/rapidgo/server/internal/adapters/webAdp"
	"github.com/vanshjangir/rapidgo/server/internal/applications"
)

func main() {
	db := postgresAdp.SetupDB()
	userRepo := postgresAdp.NewPostgresUserRepo(db)
    userService := applications.NewUserService(userRepo)

	httpHandler := webAdp.NewGinHandler(userService, nil)
	httpHandler.RegisterRoutes()
	httpHandler.Run()

	fmt.Println("Starting server")
}
