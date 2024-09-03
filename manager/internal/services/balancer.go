package services

import (
	"encoding/json"
	"errors"
	"manager/internal/docker"
	"manager/internal/utils"
)

type Application struct {
	Domain        string `json:"domain"`
	ContainerName string `json:"name"`
	Port          string `json:"port"`
	PathPrefix    string `json:"path_prefix"`
}

type BalancerState struct {
	Domains      []string `json:"domains"`
	Email        string   `json:"email"`
	Applications []Application
}

func (s *Services) GetBalancerState() (BalancerState, error) {
	body, err := s.broker.Request("controller:get:balancer-state", []byte{})
	if err != nil {
		return BalancerState{}, err
	}

	var state BalancerState
	err = json.Unmarshal(body, &state)
	return state, err
}

func (s *Services) RunHttp() error {
	body, err := s.broker.Request("balancer:launch:http", []byte{})
	if err != nil {
		return err
	}
	if !utils.IsSuccessMessage(body) {
		return errors.New("failed send launch http")
	}

	return nil
}

func (s *Services) RunHttps(domain, email string) error {
	data, err := json.Marshal(map[string]string{
		"domain": domain,
		"email":  email,
	})
	if err != nil {
		return err
	}
	body, err := s.broker.Request("balancer:launch:http", data)
	if err != nil {
		return err
	}
	if !utils.IsSuccessMessage(body) {
		return errors.New("failed send launch https")
	}

	return nil
}

func (s *Services) RunBalancer() error {
	options := docker.CreateContainerOptions{
		Image:         "aranea/balancer:latest",
		ContainerName: "balancer",
		Networks:      "aranea-network",
		Binds: []struct {
			Host      string
			Container string
		}{
			{
				Host:      "/data/aranea/",
				Container: "/data/aranea/",
			},
		},
		Ports: map[int]int{
			80:  80,
			443: 443,
		},
	}

	id, err := s.docker.CreateContainer(options)
	if err != nil {
		return err
	}

	err = s.docker.RunContainer(id)
	return err
}

func (s *Services) AddApplications(apps []Application) error {
	for _, app := range apps {
		data, err := json.Marshal(app)
		if err != nil {
			return err
		}
		body, err := s.broker.Request("balancer:add:app", data)
		if err != nil {
			return err
		}

		if !utils.IsSuccessMessage(body) {
			return errors.New("failed to add app to balancer")
		}
	}

	return nil
}

func (s *Services) AddDomains(domains []string) error {
	for _, domain := range domains {
		data, err := json.Marshal(map[string]string{
			"domain": domain,
		})
		if err != nil {
			return err
		}

		body, err := s.broker.Request("balancer:add:domain", data)
		if err != nil {
			return err
		}

		if !utils.IsSuccessMessage(body) {
			return errors.New("failed add domain to balancer")
		}
	}

	return nil
}

func (s *Services) AddControllerAndUIApplications() error {
	controller := Application{
		Domain:        "*",
		Port:          "8000",
		PathPrefix:    "/api",
		ContainerName: "controller",
	}

	data, err := json.Marshal(controller)
	if err != nil {
		return err
	}

	body, err := s.broker.Request("balancer:add:app", data)
	if err != nil {
		return err
	}

	if !utils.IsSuccessMessage(body) {
		return errors.New("failed to add controller to balancer")
	}

	ui := Application{
		Domain:        "*",
		PathPrefix:    "/",
		Port:          "3000",
		ContainerName: "aranea-ui",
	}
	data, err = json.Marshal(ui)
	if err != nil {
		return err
	}

	body, err = s.broker.Request("balancer:add:app", data)
	if err != nil {
		return err
	}

	if !utils.IsSuccessMessage(body) {
		return errors.New("failed to add ui to balancer")
	}

	return nil
}
