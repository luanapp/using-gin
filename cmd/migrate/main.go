package main

import (
	"github.com/luanapp/gin-example/cmd"
	"github.com/luanapp/gin-example/config/database"
	_ "github.com/luanapp/gin-example/pkg/env"
)

func main() {
	database.InitializeDB()
	cmd.Execute()
}
