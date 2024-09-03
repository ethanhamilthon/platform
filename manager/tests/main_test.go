package tests

import (
	"context"
	"encoding/json"
	"manager/internal/message"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SuccessMessage() []byte {
	data, _ := json.Marshal(map[string]interface{}{
		"message": "success",
	})
	return data
}

// TODO: write integration test with test containers
func TestMain(t *testing.T) {
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
	defer natsContainer.Terminate(ctx)

	// Get mapped port
	mappedPort, err := natsContainer.MappedPort(ctx, "4222")
	if err != nil {
		t.Fatal(err)
	}

	// Create message broker
	broker, err := message.New("nats://localhost:" + mappedPort.Port())
	if err != nil {
		t.Fatal(err)
	}
	defer broker.Close()

	SubTestHttpsLaunched(t, broker)
	SubTestHttpLaunched(t, broker)
}
