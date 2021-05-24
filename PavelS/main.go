package main

import (
	"github.com/SavinskiPavel/bookstore/internal/database"
	"github.com/SavinskiPavel/bookstore/internal/handler"
	"github.com/SavinskiPavel/bookstore/internal/server"
	"log"
)

func main() {
	srv := new(server.Server)
	database.InitDB()
	handler.HandleRequest()
	if err := srv.Run("8080", nil); err != nil {
		log.Fatal(err)
	}
}
