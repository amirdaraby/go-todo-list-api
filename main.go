package main

import (
	"log"

	"github.com/amirdaraby/go-todo-list-api/internal/config"
	"github.com/amirdaraby/go-todo-list-api/internal/db"
)

func main() {
	config, err := config.Init()

	if err != nil {
		log.Fatalf("config init failed, err: %s", err)
	}

	_ , err = db.Init(*config)

	if err != nil {
		log.Fatalf("db init failed, err: %s", err)
	}

}
