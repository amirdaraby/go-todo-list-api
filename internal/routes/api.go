package routes

import (
	"github.com/amirdaraby/go-todo-list-api/internal/handlers"
	"github.com/amirdaraby/go-todo-list-api/internal/middleware"
	"github.com/gorilla/mux"
)

func Init() *mux.Router {

	// main router
	r := mux.NewRouter()

	// api router
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.JsonRequestAndResponse)

	// api version one router
	apiV1 := api.PathPrefix("/v1").Subrouter()

	authApis := apiV1.PathPrefix("/auth").Subrouter()
	authApis.Use(middleware.Guest)

	authApis.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	authApis.HandleFunc("/login", handlers.LoginUser).Methods("POST")

	userApis := apiV1.PathPrefix("/user").Subrouter()
	userApis.Use(middleware.Authenticate)

	userApis.HandleFunc("", handlers.ShowUser).Methods("GET")
	userApis.HandleFunc("", handlers.UpdateUser).Methods("PATCH")

	taskApis := userApis.PathPrefix("/tasks").Subrouter()
	taskApis.HandleFunc("", handlers.IndexTask).Methods("GET")
	taskApis.HandleFunc("/undone", handlers.IndexUnDoneTasks).Methods("GET")
	taskApis.HandleFunc("/done", handlers.IndexDoneTasks).Methods("GET")
	taskApis.HandleFunc("", handlers.StoreTask).Methods("POST")
	taskApis.HandleFunc("/{id}", handlers.ShowTask).Methods("GET")
	taskApis.HandleFunc("/{id}", handlers.DeleteTask).Methods("DELETE")
	taskApis.HandleFunc("/{id}", handlers.UpdateTask).Methods("PATCH")
	taskApis.HandleFunc("/{id}/toggle-done-status", handlers.ToggleDoneTask).Methods("POST")

	return r
}
