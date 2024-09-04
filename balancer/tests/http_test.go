package tests

import (
	"balancer/config"
	"balancer/internal/message"
	"balancer/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	// "time"
)

func HttpLaunchTest(broker *message.Message, t *testing.T) {
	t.Log("Test http")
	ok := isPortAvailable("localhost", config.HttpPort)
	if !ok {
		t.Fatal("port 8080 not available")
	}

	t.Log(config.HttpTopic)
	msg, err := broker.Request(config.HttpTopic, []byte{})
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

	// time.Sleep(20 * time.Second)
	// Check if http was launched
	checkHttpServer(t, "http://localhost:"+config.HttpPort)
}

func isPortAvailable(host string, port string) bool {
	address := fmt.Sprintf("%s:%s", host, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		// Если произошла ошибка, это значит, что порт уже занят
		return false
	}
	// Если удалось открыть сокет, закрываем его и возвращаем true
	defer listener.Close()
	return true
}

func checkHttpServer(t *testing.T, url string) {
	resp, err := http.Get(url)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	// Проверяем статус-код
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Ожидался статус 404, но получен %d", resp.StatusCode)
	}

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	// Проверяем содержимое тела
	expectedBody := "not found"
	if string(body) != expectedBody {
		t.Errorf("Ожидалось тело ответа '%s', но получено '%s'", expectedBody, string(body))
	}
}
