package tests

import (
	"context"
	"encoding/json"
	"errors"
	"manager/internal/message"
	"manager/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SuccessMessage() []byte {
	data, _ := json.Marshal(map[string]interface{}{
		"message": "success",
	})
	return data
}

func AddBrokerSubscriberIfHttps(t *testing.T, broker *message.MessageBroker) {
	broker.Response("controller:get:balancer-state", func(data []byte) ([]byte, error) {
		data, err := json.Marshal(map[string]interface{}{
			"domains": []string{"app.com", "admin.app.com", "app2.com"},
			"email":   "test@app.com",
			"applications": []services.Application{
				{
					Domain:        "app.com",
					ContainerName: "test",
					Port:          "80",
					PathPrefix:    "",
				},
				{
					Domain:        "app2.com",
					ContainerName: "test",
					Port:          "80",
					PathPrefix:    "",
				},
				{
					Domain:        "admin.app.com",
					ContainerName: "test",
					Port:          "80",
					PathPrefix:    "",
				},
			},
		})
		if err != nil {
			assert.NoError(t, err)
			return []byte{}, err
		}
		return data, nil
	})

	broker.Response("balancer:launch:https", func(data []byte) ([]byte, error) {
		var body map[string]interface{}
		err := json.Unmarshal(data, &body)
		if err != nil {
			assert.NoError(t, err)
			return []byte{}, err
		}
		if body["domain"] == "app.com" && body["email"] == "test@app.com" {
			return SuccessMessage(), nil
		}

		assert.Error(t, errors.New("failed send launch https"))
		return []byte{}, errors.New("failed send launch https")
	})

	broker.Response("balancer:launch:http", func(data []byte) ([]byte, error) {
		return SuccessMessage(), nil
	})

	broker.Response("balancer:add:app", func(data []byte) ([]byte, error) {
		var app services.Application
		err := json.Unmarshal(data, &app)
		if err != nil {
			assert.NoError(t, err)
			return []byte{}, err
		}
		return SuccessMessage(), nil
	})

	broker.Response("balancer:add:domain", func(data []byte) ([]byte, error) {
		var domain struct {
			Domain string `json:"domain"`
		}
		err := json.Unmarshal(data, &domain)
		if err != nil {
			assert.NoError(t, err)
			return []byte{}, err
		}
		return SuccessMessage(), nil
	})
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
		assert.NoError(t, err)
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
		assert.NoError(t, err)
	}
	defer broker.Close()

	AddBrokerSubscriberIfHttps(t, broker)
	//  Create  new service
	srv := services.New(broker, nil)

	// Start platform
	err = srv.StartPlatform()
	if err != nil {
		assert.NoError(t, err)
	}
}
