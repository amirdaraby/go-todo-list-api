package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.AppPort),
		Handler: router,
	}

	gracefulShutdownCh := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		cause := <-sigint

		log.Printf("got %s signal, shutting down...", cause.String())

		server.Shutdown(context.Background())

		close(gracefulShutdownCh)
	}()

	err = server.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}

	<-gracefulShutdownCh
}
