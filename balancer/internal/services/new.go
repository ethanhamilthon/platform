package services

import (
	"balancer/config"
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
	Email   string `json:"email"`
	Domains string `json:"domain"`
}

func (s *Service) LaunchHttp(data []byte) ([]byte, error) {
	// Routers
	router := mux.NewRouter()
	router.HandleFunc("/.well-known/acme-challenge/", func(w http.ResponseWriter, r *http.Request) {
		// Handle HTTP-01 challenges
		certmagic.DefaultACME.HandleHTTPChallenge(w, r)
	})
	router.PathPrefix("/").HandlerFunc(s.Proxy)
	server := &http.Server{
		Addr:    ":" + config.HttpPort,
		Handler: router,
	}

	// Serve listener
	go func() {
		log.Printf("Http server started on port %s\n", config.HttpPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v\n", err)
		}
	}()

	return utils.Success()
}

func (s *Service) LaunchHttps(data []byte) ([]byte, error) {
	var body LaunchData
	err := json.Unmarshal(data, &body)
	if err != nil {
		return []byte{}, err
	}
	// Start server
	router := mux.NewRouter()
	router.PathPrefix("/").HandlerFunc(s.Proxy)
	if config.Mode != "dev" {
		// Configure CertMagic
		certmagic.DefaultACME.AltHTTPPort = -1
		certmagic.DefaultACME.DisableHTTPChallenge = true
		certmagic.DefaultACME.Email = body.Email
		certmagic.Default.Storage = &certmagic.FileStorage{Path: "/certs/storage"}

		go func() {
			if err := certmagic.HTTPS([]string{body.Domains}, router); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Could not start server: %v\n", err)
			}
		}()

	} else {

		cert, key, err := utils.GenerateSelfSignedCert()
		if err != nil {
			return []byte{}, err
		}
		go func() {
			if err := http.ListenAndServeTLS(":"+config.HttpsPort, cert, key, router); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Could not start server: %v\n", err)
			}
		}()

	}
	return utils.Success()
}
