package main

import (
	"fmt"
	"os"

	"github.com/lfmexi/gochat/handlers"
	"github.com/lfmexi/gochat/server"
)

var host = "0.0.0.0"
var port = "8888"

func init() {
	if envhost := os.Getenv("GOCHAT_HOST"); envhost != "" {
		host = envhost
	}
	if envport := os.Getenv("GOCHAT_PORT"); envport != "" {
		port = envport
	}
}

func main() {
	server := server.NewServer(fmt.Sprintf("%s:%s", host, port))
	handler := handlers.NewChatHandler()
	server.SetHandler(handler)
	server.Listen()
}
