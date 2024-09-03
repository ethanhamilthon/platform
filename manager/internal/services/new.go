package services

import (
	"fmt"
	"manager/internal/docker"
	"manager/internal/message"
)

type Services struct {
	broker *message.MessageBroker
	docker *docker.DockerSerive
}

func New(broker *message.MessageBroker, docker_client *docker.DockerSerive) *Services {
	return &Services{broker: broker, docker: docker_client}
}

func (s *Services) StartPlatform() error {
	// Check the ports
	// err := s.CheckPorts([]int{443, 80})
	// if err != nil {
	// 	return err
	// }
	//
	// // Run Controller, Adapter, UI, Balancer
	// err = s.RunAdapter()
	// if err != nil {
	// 	return err
	// }
	//
	// err = s.RunUI()
	// if err != nil {
	// 	return err
	// }
	//
	// err = s.RunController()
	// if err != nil {
	// 	return err
	// }
	//
	// err = s.RunBalancer()
	// if err != nil {
	// 	return err
	// }

	// Get Balancer state
	state, err := s.GetBalancerState()
	if err != nil {
		return err
	}

	// Check if we can run https right away
	isHttpsAvailable := false
	if len(state.Domains) > 0 {
		isHttpsAvailable = true
	}

	// if we can run https ()
	if isHttpsAvailable {
		err = s.RunHttp()
		if err != nil {
			return err
		}

		err = s.RunHttps(state.Domains[0], state.Email)
		if err != nil {
			return err
		}

		state.Domains = state.Domains[1:]
		err = s.AddDomains(state.Domains)
		if err != nil {
			return err
		}

		// Add all applications (include platform controller api and ui)
		err = s.AddApplications(state.Applications)
		if err != nil {
			return err
		}
	} else {
		// If we can not run https (since it is first launch, and use have not entered domain, email yet)
		err = s.RunHttp()
		if err != nil {
			return err
		}

		err = s.AddControllerAndUIApplications()
		if err != nil {
			return err
		}
	}

	fmt.Println("All services successfully runned")
	return nil
}
