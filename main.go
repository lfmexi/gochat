package main

import (
	"fmt"
	"log"
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
	exit := make(chan os.Signal, 1)
	serverShutdown := make(chan bool)

	server := server.NewServer(fmt.Sprintf("%s:%s", host, port), serverShutdown)
	handler := handlers.NewChatHandler()
	server.SetHandler(handler)
	server.Listen()

	sig := <-exit

	log.Printf("Exiting with signal %s", sig)

	serverShutdown <- true
}
