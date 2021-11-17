package main

import (
	"github.com/luanapp/gin-example/config/database"
	"github.com/luanapp/gin-example/pkg/server"
)

// @title Natural History Museum API documentation
// @version 1.0.0
// @host localhost:8080
// @BasePath /
func main() {
	database.InitializeDB()
	server.NewServer().Start()
}
