package tests

import (
	"context"
	"controller/cmd"
	"controller/internal/config"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	token = ""
)

func TestMain(t *testing.T) {
	// Setup
	dbPath := "./auth.db"
	removeFileIfExists(dbPath)
	defer removeFileIfExists(dbPath)
	createEmptyFile(dbPath)
	config.Init()
	config.DbPath = dbPath

	// Run nats
	mappedPort, closer := runNats(t)
	defer closer()
	config.NatsUrl = "nats://localhost:" + mappedPort
	// Run server
	go cmd.Start()

	// Wait for server
	waitForTCPServer("localhost:8000", t)

	// Run tests
	t.Run("AuthTest", AuthTest)
	t.Run("CreateAppTest", createAppTest)
}
func waitForTCPServer(address string, t *testing.T) {
	tries := 200
	for {
		if tries == 0 {
			t.Errorf("timeout waiting for tcp server")
			break
		}
		conn, err := net.DialTimeout("tcp", address, 1*time.Second)
		if err == nil {
			conn.Close()
			break
		}
		time.Sleep(100 * time.Millisecond) // небольшая пауза между попытками
		tries--
	}
}
func createEmptyFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	fmt.Println("File created:", filename)
	return nil
}

func removeFileIfExists(filename string) error {
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		err := os.Remove(filename)
		if err != nil {
			return fmt.Errorf("error removing file: %w", err)
		}
		fmt.Println("File removed:", filename)
	} else {
		fmt.Println("File does not exist:", filename)
	}
	return nil
}

func runNats(t *testing.T) (string, func()) {
	ctx := context.Background()

	// Container  config
	req := testcontainers.ContainerRequest{
		Image:        "nats:latest",
		ExposedPorts: []string{"4222/tcp"},
		WaitingFor:   wait.ForListeningPort("4222/tcp"),
		Cmd:          []string{"-js"},
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
