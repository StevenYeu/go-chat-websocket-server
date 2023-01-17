package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/StevenYeu/go-chat-websocket-server/server"
	"nhooyr.io/websocket"
)

func main() {
	l, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Print("network error")
	}
	s := &http.Server{
		Handler: server.ChatServer{
			Clients: make(map[int]*websocket.Conn),
		},
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	errc := make(chan error, 1)
	fmt.Printf("Starting Chat Server on localhost:8080\n")
	go func() {
		errc <- s.Serve(l)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	s.Shutdown(ctx)
}
