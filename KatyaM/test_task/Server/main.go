// +build !solution
package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	handler2 "test_task/Server/handler"
	service3 "test_task/Server/service"
	"time"
)

var (
	POSTGRES_USER     = "user"
	POSTGRES_PASSWORD = "123"
	POSTGRES_DB       = "postgres"
)

func main() {
	addr := ":8080"
	dbUser, dbPassword, dbName := POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB
	listener, _ := net.Listen("tcp", addr)
	l := zap.NewExample()
	service := service3.BookService{L: l}
	db, _ := service.Initialize(dbUser, dbPassword, dbName)
	httpHandler := handler2.NewBookHandler(l, service, db)
	mux := http.NewServeMux()
	httpHandler.Register(mux)
	server := http.Server{Handler: mux}
	go func() {
		server.Serve(listener)
	}()
	defer Stop(&server)
	log.Printf("Started server on %s", addr)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(fmt.Sprint(<-ch))
	log.Println("Stopping API server.")
}

func Stop(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Could not shut down server correctly: %v\n", err)
		os.Exit(1)
	}
}
