package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

type Publisher interface {
	Publish(topic string, data []byte) error
	Request(topic string, data []byte) ([]byte, error)
}

type Services struct {
	Broker Publisher
}

func New() *Services {
	return &Services{}
}

func (s *Services) Load() {
	// Check the ports
	ports := [4]int{8000, 80, 443, 3000}
	for i := 0; i < len(ports); i++ {
		if !isPortAvailable(ports[i]) {
			log.Fatalf("port %v is not available", ports[i])
		}
	}

	// Run Services

	// Run controller rest server
	err := s.Broker.Publish("controller:run:server", []byte{})
	if err != nil {
		log.Fatalln("error publish to controller")
	}

	// Create ui container
	body, err := s.Broker.Request("adapter:create:container", []byte{}) // Todo: add real params for launching ui
	if err != nil {
		log.Fatalln("error publish create ui container")
	}
	type IdBody struct {
		ID string `json:"id"`
	}
	var id_body IdBody
	err = json.Unmarshal(body, &id_body)
	if err != nil {
		log.Fatalln("error unmarshal json for id body")
	}

	// Run ui container
	body, err = json.Marshal(map[string]string{
		"id": id_body.ID,
	})
	err = s.Broker.Publish("adapter:run:container", body)
	if err != nil {
		log.Fatalln("error publish run container")
	}

	// Launch http
	err = s.Broker.Publish("balancer:launch:http", []byte{})
	if err != nil {
		log.Fatalln("error to publish launch http")
	}

	// Get domains from controller
	type DomainBody struct {
		Domains []string `json:"domains"`
	}

	body, err = s.Broker.Request("controller:get:domains", []byte{})
	if err != nil {
		log.Fatalln("error request controller to get domains")
	}

	var domain_body DomainBody
	err = json.Unmarshal(body, &domain_body)
	if err != nil {
		log.Fatalln("error request unmarshal json")
	}

	// If there are domains
	if len(domain_body.Domains) != 0 {
		// Get email
		type EmailBody struct {
			Email string `json:"email"`
		}

		body, err = s.Broker.Request("controller:get:email", []byte{})
		if err != nil {
			log.Fatalln("error request controller to get email")
		}

		var email_body EmailBody
		err = json.Unmarshal(body, &email_body)
		if err != nil {
			log.Fatalln("error unmarshal email_body")
		}

		// Run https
		body, err = json.Marshal(map[string]interface{}{
			"email":   email_body.Email,
			"domains": domain_body.Domains,
		})

		if err != nil {
			log.Fatalln("error marshal body for launch https")
		}

		err = s.Broker.Publish("balancer:launch:https", body)
		if err != nil {
			log.Fatalln("error launch https")
		}

		// Todo: launch all apps in system
	}

	// Todo Add containers
}

func isPortAvailable(port int) bool {
	addr := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	ln.Close()
	return true
}
