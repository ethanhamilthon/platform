package tests

import (
	"balancer/config"
	"balancer/internal/message"
	"balancer/internal/services"
	"balancer/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func HttpProxyTest(t *testing.T, broker *message.Message) {
	t.Log("Test http proxy")
	// Add apps
	addNewApp(t, broker, services.AddAppBody{
		Domain:    "localhost",
		Path:      "/",
		Container: "test1",
		Port:      "3000",
	})
	addNewApp(t, broker, services.AddAppBody{
		Domain:    "localhost",
		Path:      "/test",
		Container: "test2",
		Port:      "3001",
	})

	// Do requests
	err := doHttpProxyRequest("http://localhost/", "http://test1:3000", 200)
	if err != nil {
		t.Fatal(err)
	}
	err = doHttpProxyRequest("http://localhost/test", "http://test2:3001", 200)
	if err != nil {
		t.Fatal(err)
	}
}

func addNewApp(t *testing.T, broker *message.Message, app services.AddAppBody) {
	t.Log(app)
	body, err := json.Marshal(app)
	if err != nil {
		t.Fatal(err)
	}
	msg, err := broker.Request(config.AddAppTopic, body)
	if err != nil {
		t.Fatal(err)
	}
	data := new(utils.Message)
	err = json.Unmarshal(msg.Data, data)
	if err != nil {
		t.Fatal(err)
	}
	if data.M != "success" {
		t.Fatal("message not success")
	}
}

func doHttpProxyRequest(url string, expectBody string, expectCode int) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != expectCode {
		return fmt.Errorf("Expected status code %d, but got %d", expectCode, resp.StatusCode)
	}

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Check the response body
	if string(body) != expectBody {
		return fmt.Errorf("Expected body '%s', but got '%s'", expectBody, string(body))
	}

	return nil
}
