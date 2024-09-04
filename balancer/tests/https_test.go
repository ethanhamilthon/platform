package tests

import (
	"balancer/config"
	"balancer/internal/message"
	"balancer/internal/utils"
	"encoding/json"
	"testing"
)

func HttpsLaunchTest(t *testing.T, broker *message.Message) {
	t.Log("Test https")
	ok := isPortAvailable("localhost", config.HttpsPort)
	if !ok {
		t.Fatal("port 8443 not available")
	}

	body, err := json.Marshal(map[string]string{
		"domain": "localhost",
		"email":  "app@app.com",
	})

	t.Log(config.HttpsTopic)
	msg, err := broker.Request(config.HttpsTopic, body)
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
