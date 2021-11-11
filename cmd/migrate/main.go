package main

import (
	"github.com/luanapp/gin-example/cmd"
	"github.com/luanapp/gin-example/config/database"
)

func main() {
	database.InitializeDB()
	cmd.Execute()
}
