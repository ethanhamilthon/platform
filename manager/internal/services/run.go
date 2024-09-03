package services

import "manager/internal/docker"

func (s *Services) RunUI() error {
	options := docker.CreateContainerOptions{
		Image:         "aranea/ui:latest",
		ContainerName: "aranea-ui",
		Networks:      "aranea-network",
	}

	id, err := s.docker.CreateContainer(options)
	if err != nil {
		return err
	}

	err = s.docker.RunContainer(id)
	return err
}

func (s *Services) RunController() error {
	options := docker.CreateContainerOptions{
		Image:         "aranea/controller:latest",
		ContainerName: "controller",
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
	}

	id, err := s.docker.CreateContainer(options)
	if err != nil {
		return err
	}

	err = s.docker.RunContainer(id)
	return err
}

func (s *Services) RunAdapter() error {
	options := docker.CreateContainerOptions{
		Image:         "aranea/adapter:latest",
		ContainerName: "adapter",
		Networks:      "aranea-network",
		Binds: []struct {
			Host      string
			Container string
		}{
			{
				Host:      "/data/aranea/",
				Container: "/data/aranea/",
			},
			{
				Host:      "/var/run/docker.sock",
				Container: "/var/run/docker.sock",
			},
		},
	}

	id, err := s.docker.CreateContainer(options)
	if err != nil {
		return err
	}

	err = s.docker.RunContainer(id)
	return err
}
