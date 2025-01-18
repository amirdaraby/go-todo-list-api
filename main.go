package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/amirdaraby/go-todo-list-api/internal/config"
	"github.com/amirdaraby/go-todo-list-api/internal/db"
	"github.com/amirdaraby/go-todo-list-api/internal/routes"
)

func main() {
	config, err := config.Init()

	if err != nil {
		log.Fatalf("config init failed, err: %s", err)
	}

	_, err = db.Init(*config)

	if err != nil {
		log.Fatalf("db init failed, err: %s", err)
	}

	router := routes.Init()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.AppPort), router))
}
