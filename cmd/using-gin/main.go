package main

import (
	"github.com/luanapp/gin-example/config/database"
	"github.com/luanapp/gin-example/pkg/server"
)

func main() {
	database.InitializeDB()
	server.NewServer().Start()
}
