package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Init() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	return r
}
