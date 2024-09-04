package tests

import (
	"balancer/cmd"
	"balancer/config"
	"balancer/internal/message"
	"context"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestMain(t *testing.T) {
	// Initial
	mappedPort, closer := runNats(t)
	defer closer()
	loadConfigs("nats://localhost:" + mappedPort)
	cmd.Start()

	// Create client broker
	broker := message.New("nats://localhost:" + mappedPort)
	HttpLaunchTest(broker, t)
	HttpProxyTest(t, broker)
	HttpsLaunchTest(t, broker)
}

func loadConfigs(nats_url string) {
	config.Init()
	config.Mode = "dev"
	config.NatsUrl = nats_url
	// config.HttpPort = "8080"
	// config.HttpsPort = "8443"
	config.SetTesting()
}

func runNats(t *testing.T) (string, func()) {
	ctx := context.Background()

	// Container  config
	req := testcontainers.ContainerRequest{
		Image:        "nats:latest",
		ExposedPorts: []string{"4222/tcp"},
		WaitingFor:   wait.ForListeningPort("4222/tcp"),
	}

	natsContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}
	closer := func() {
		natsContainer.Terminate(ctx)
	}

	// Get mapped port
	mappedPort, err := natsContainer.MappedPort(ctx, "4222")
	if err != nil {
		t.Fatal(err)
	}

	return mappedPort.Port(), closer
}
