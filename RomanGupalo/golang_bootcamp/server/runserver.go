package main

import (
	"log"
	"net/http"
	"time"

	"github.com/AA55hex/golang_bootcamp/server/config"
	"github.com/AA55hex/golang_bootcamp/server/db/connection"
	"github.com/AA55hex/golang_bootcamp/server/handlers"
	"github.com/gorilla/mux"
	"github.com/upper/db/v4/adapter/mysql"
)

func main() {
	err := config.LoadConfigs("configs.env")
	if err != nil {
		log.Fatal("configs.env file not found: ", err)
	}
	// creating database session
	db_settings := mysql.ConnectionURL{
		Database: config.MySQL.Database,
		Host:     config.MySQL.Host,
		User:     config.MySQL.User,
		Password: config.MySQL.Password,
	}
	_, err = connection.OpenSession(&db_settings, config.MySQL.ConnectionTryCount)
	if err != nil {
		log.Fatal("Session not created: ", err)
	}
	defer connection.GetSession().Close()

	err = connection.TryMigrate(config.MySQL.MigrationSource)
	if err != nil {
		log.Fatal("Migration error: ", err)
	}

	// init router
	log.Println("Creating router")
	router := mux.NewRouter()
	router.HandleFunc("/books/{id:[0-9]+}", handlers.GetBookByIDHandler).Methods("GET")
	router.HandleFunc("/books/{id:[0-9]+}", handlers.UpdateBookHandler).Methods("PUT")
	router.HandleFunc("/books/{id:[0-9]+}", handlers.DeleteBookHandler).Methods("DELETE")
	router.HandleFunc("/books", handlers.GetBooksByFilterHandler).Methods("GET")
	router.HandleFunc("/books/new", handlers.CreateBookHandler).Methods("POST")

	// listen & serve
	log.Println("Creating server")
	server := http.Server{
		Handler:      router,
		Addr:         config.Server.Address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Server created.")
	log.Println("Listening started on: ", server.Addr)
	log.Fatal(server.ListenAndServe())
}
