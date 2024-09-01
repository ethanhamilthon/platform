package services

import (
	"balancer/internal/utils"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/caddyserver/certmagic"
	"github.com/gorilla/mux"
)

type AppType struct {
	Port          string `json:"port"`
	Path          string `json:"path"`
	ContainerName string `json:"container_name"`
}

type DomainType struct {
	Domain string    `json:"domain"`
	Apps   []AppType `json:"apps"`
}

type Service struct {
	mu              *sync.RWMutex
	isProxyLaunched bool
	domains         []DomainType
}

func New() *Service {
	return &Service{
		mu:              new(sync.RWMutex),
		isProxyLaunched: false,
		domains:         []DomainType{},
	}
}

type LaunchData struct {
	Email   string   `json:"email"`
	Domains []string `json:"domains"`
}

func (s *Service) LaunchHttp(data []byte) ([]byte, error) {
	// Routers
	router := mux.NewRouter()
	router.PathPrefix("/").HandlerFunc(s.Proxy)
	server := &http.Server{
		Addr:    ":80",
		Handler: router,
	}

	// Serve listener
	go log.Fatal(server.ListenAndServe())

	return utils.Success()
}

func (s *Service) LaunchHttps(data []byte) ([]byte, error) {
	var body LaunchData
	err := json.Unmarshal(data, &body)
	if err != nil {
		return []byte{}, err
	}
	// Configure CertMagic
	certmagic.DefaultACME.AltHTTPPort = -1
	certmagic.DefaultACME.DisableHTTPChallenge = true
	certmagic.DefaultACME.Email = body.Email
	certmagic.Default.Storage = &certmagic.FileStorage{Path: "/certs/storage"}

	// Start server
	router := mux.NewRouter()
	router.PathPrefix("/").HandlerFunc(s.Proxy)
	go log.Fatal(certmagic.HTTPS(body.Domains, router))

	return utils.Success()
}
