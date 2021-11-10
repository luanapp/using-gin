package main

import (
	"luana.com/gin-example/config/database"
	"luana.com/gin-example/pkg/server"
)

func main() {
	database.InitializeDB()
	server.NewServer().Start()
}
