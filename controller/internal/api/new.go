package api

import (
	"controller/internal/services/application"
	"controller/internal/services/auth"
	"controller/internal/services/domain"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	auth   *auth.AuthService
	domain *domain.DomainService
	app    *application.ApplicationService
}

func New(auth *auth.AuthService, dm *domain.DomainService, app *application.ApplicationService) *ApiServer {
	return &ApiServer{auth: auth, domain: dm, app: app}
}

func (a *ApiServer) Serve() {
	router := mux.NewRouter()
	// Api router
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/auth/login", a.LoginHandler).Methods("POST")
	api.HandleFunc("/auth/register", a.RegisterHandler).Methods("POST")
	api.HandleFunc("/auth/me", a.MeHandler)
	api.HandleFunc("/domain/list", a.Protected(a.DomainListHandler)).Methods("GET")
	api.HandleFunc("/domain/add", a.Protected(a.DomainAddHandler)).Methods("POST")
	api.HandleFunc("/github/pull", a.Protected(a.PublicGithubPullHandler)).Methods("POST")
	api.HandleFunc("/repo/detect", a.Protected(a.DetectRepoHandler)).Methods("POST")
	api.HandleFunc("/app/create", a.Protected(a.CreateApp)).Methods("POST")
	// Create server
	server := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	log.Fatalln(server.ListenAndServe())
}

func (a *ApiServer) Protected(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		err := a.auth.VerifyJWT(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}
