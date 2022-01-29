package main

import (
	"github.com/luanapp/gin-example/cmd"
	_ "github.com/luanapp/gin-example/pkg/env"
)

func main() {
	cmd.Execute()
}
