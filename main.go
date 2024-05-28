package main

import (
	"github.com/Surya-7890/gokey/server"
)

func main() {
	newServer := server.NewServer(":7000")
	newServer.StartServer()
}
