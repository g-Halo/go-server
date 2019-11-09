package main

import (
	"github.com/g-Halo/go-server/logger"
)

func init() {
	logger.InitLogger("./auth.log", "debug")
}

func main() {
	StartRpc()
}