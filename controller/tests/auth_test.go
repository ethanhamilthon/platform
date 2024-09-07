package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func AuthTest(t *testing.T) {
	t.Run("RegisterTest", registerTest)
	t.Run("LoginTest", loginTest)
}

func registerTest(t *testing.T) {
	// Register
	body, err := doHttpRequest("http://localhost:8000/api/auth/register", "POST",
		map[string]string{"username": "testuser", "password": "password123", "confirm_password": "password123"})

	if err != nil {
		t.Errorf("error registering user: %v", err)
	}

	if len(body) == 0 {
		t.Errorf("error registering user: empty body")
	}

	data := make(map[string]string, 0)
	err = json.Unmarshal(body, &data)
	if err != nil {
		t.Errorf("error registering user: %v, body: %v", err, string(body))
	}

	if data["token"] == "" {
		t.Errorf("error registering user: empty token")
	}

	token = data["token"]

	body, err = doHttpRequest("http://localhost:8000/api/auth/me", "GET", map[string]string{"token": token})

	if err != nil {
		t.Errorf("error getting user: %v", err)
	}

	if len(body) == 0 {
		t.Errorf("error getting user: empty body")
	}

	data = make(map[string]string, 0)
	err = json.Unmarshal(body, &data)
	if err != nil {
		t.Errorf("error getting user: %v, body: %v", err, string(body))
	}

	if data["username"] != "testuser" {
		t.Errorf("error getting user: wrong username")
	}

	t.Log("RegisterTest: PASS")
}

func loginTest(t *testing.T) {
	// Login
	body, err := doHttpRequest("http://localhost:8000/api/auth/login", "POST",
		map[string]string{"username": "testuser", "password": "password123"})

	if err != nil {
		t.Errorf("error logging in: %v", err)
	}

	if len(body) == 0 {
		t.Errorf("error logging in: empty body")
	}

	data := make(map[string]string, 0)
	err = json.Unmarshal(body, &data)
	if err != nil {
		t.Errorf("error logging in: %v, body: %v", err, string(body))
	}

	if data["token"] == "" {
		t.Errorf("error logging in: empty token")
	}

	token = data["token"]

	body, err = doHttpRequest("http://localhost:8000/api/auth/me", "GET", map[string]string{"token": token})

	if err != nil {
		t.Errorf("error getting user: %v", err)
	}

	if len(body) == 0 {
		t.Errorf("error getting user: empty body")
	}

	data = make(map[string]string, 0)
	err = json.Unmarshal(body, &data)
	if err != nil {
		t.Errorf("error getting user: %v, body: %v", err, string(body))
	}

	if data["username"] != "testuser" {
		t.Errorf("error getting user: wrong username")
	}

	t.Log("LoginTest: PASS")
}

func doHttpRequest(url string, method string, data any) ([]byte, error) {
	// Turn into json
	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshalling body: %w", err)
	}
	// Setup
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	// Execute
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	defer res.Body.Close()
	body, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	return body, nil
}
