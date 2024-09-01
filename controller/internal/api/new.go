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

	// Create routes

	// Create server
	server := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	log.Fatalln(server.ListenAndServe())
}
