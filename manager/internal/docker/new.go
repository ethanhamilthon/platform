package docker

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type DockerSerive struct {
	client *client.Client
}

func New() (*DockerSerive, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return &DockerSerive{}, errors.New("failed run docker client")
	}
	return &DockerSerive{client: cli}, nil
}

type CreateContainerOptions struct {
	Image         string
	ContainerName string
	Ports         map[int]int
	Binds         []struct {
		Host      string
		Container string
	}
	Envs     []string
	Networks string
}

func (d *DockerSerive) CreateContainer(options CreateContainerOptions) (string, error) {
	// Convert ports
	exposedPorts := make(nat.PortSet)
	bindPorts := make(nat.PortMap)
	for k, v := range options.Ports {
		exposedPorts[nat.Port(strconv.Itoa(v)+"/tcp")] = struct{}{}
		bindPorts[nat.Port(strconv.Itoa(v))] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: strconv.Itoa(k),
			},
		}
	}

	// Convert binds
	binds := make([]string, 0)
	for _, v := range options.Binds {
		b := fmt.Sprintf("%v:%v", v.Host, v.Container)
		binds = append(binds, b)
	}
	// Create container
	ctx := context.Background()
	container, err := d.client.ContainerCreate(ctx, &container.Config{
		Image:        options.Image,
		Env:          options.Envs,
		ExposedPorts: exposedPorts,
	}, &container.HostConfig{
		Binds:        binds,
		PortBindings: bindPorts,
		NetworkMode:  container.NetworkMode(options.Networks),
	}, nil, nil, options.ContainerName)
	if err != nil {
		return "", err
	}

	return container.ID, nil
}

func (d *DockerSerive) RunContainer(id string) error {
	ctx := context.Background()
	err := d.client.ContainerStart(ctx, id, container.StartOptions{})
	return err
}

func (d *DockerSerive) Inspect(id string) (types.ContainerJSON, error) {
  return d.client.ContainerInspect(context.Background(), id) 
}

func (d *DockerSerive) RemoveContainer(id string) error {
    ctx := context.Background()

    // Остановка контейнера
    if err := d.client.ContainerStop(ctx, id, container.StopOptions{}); err != nil {
        return fmt.Errorf("failed to stop container: %w", err)
    }

    // Удаление контейнера
    if err := d.client.ContainerRemove(ctx, id, container.RemoveOptions{Force: true}); err != nil {
        return fmt.Errorf("failed to remove container: %w", err)
    }

    return nil
}
