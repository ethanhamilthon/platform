package services

import (
	"fmt"
	"manager/internal/docker"
	"time"
)

func (s *Services) CheckPorts(ports []int) error {
  // Run test containers to test if 80, 443 ports are available
  fmt.Println("Start checking 80, 433 ports")

  // Container options
  options := docker.CreateContainerOptions{
    Image: "aranea/portchecker:latest",
    ContainerName: "portchecker",
    Ports: map[int]int{
      80:80,
      443:443,
    },
    Binds: []struct{Host string; Container string}{
      struct{Host string; Container string}{Host: "/data/aranea/", Container: "/data/aranea/"},
    },
  }

  id, err := s.docker.CreateContainer(options)
  if err != nil {
    return err
  }
  
  err = s.docker.RunContainer(id)
  if err != nil {
    return err
  }

  time.Sleep(time.Second)
  err = s.docker.RemoveContainer(id)
  return err
}
