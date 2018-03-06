package main

import (
	"ICityDataEngine/server"
)

func main() {
	go server.Start()
	server.StartWebSocketServer()
}
