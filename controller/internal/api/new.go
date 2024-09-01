package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiServer struct {
}

func New() *ApiServer {
	return &ApiServer{}
}

func (a *ApiServer) Serve() {
	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()
	// Create routes
	api.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {})
	api.HandleFunc("/auth/register", func(w http.ResponseWriter, r *http.Request) {})
	api.HandleFunc("/auth/me", func(w http.ResponseWriter, r *http.Request) {})
	// Create server
	server := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	log.Fatalln(server.ListenAndServe())
}
