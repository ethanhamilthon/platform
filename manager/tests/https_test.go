package tests

import (
	"encoding/json"
	"errors"
	"manager/internal/message"
	"manager/internal/services"
	"testing"
)

func AddBrokerSubscriberIfHttps(t *testing.T, broker *message.MessageBroker) func() error {
	domains := make([]string, 0)
	email := ""
	apps := make([]services.Application, 0)
	httpsLaunched := false
	httpLaunched := false

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
			t.Error(err)
			return []byte{}, err
		}
		return data, nil
	})

	broker.Response("balancer:launch:https", func(data []byte) ([]byte, error) {
		t.Log("https listening")
		var body map[string]interface{}
		err := json.Unmarshal(data, &body)
		if err != nil {
			return []byte{}, err
		}
		if body["domain"] == "app.com" && body["email"] == "test@app.com" {
			domains = append(domains, body["domain"].(string))
			email = body["email"].(string)
			httpsLaunched = true
			return SuccessMessage(), nil
		}

		return []byte{}, errors.New("failed send launch https")
	})

	broker.Response("balancer:launch:http", func(data []byte) ([]byte, error) {
		httpLaunched = true
		t.Log("http listening")
		return SuccessMessage(), nil
	})

	broker.Response("balancer:add:app", func(data []byte) ([]byte, error) {
		var app services.Application
		err := json.Unmarshal(data, &app)
		if err != nil {
			return []byte{}, err
		}
		apps = append(apps, app)
		return SuccessMessage(), nil
	})

	broker.Response("balancer:add:domain", func(data []byte) ([]byte, error) {
		var domain struct {
			Domain string `json:"domain"`
		}
		err := json.Unmarshal(data, &domain)
		if err != nil {
			return []byte{}, err
		}
		domains = append(domains, domain.Domain)
		return SuccessMessage(), nil
	})

	checkData := func() error {
		if !httpLaunched {
			return errors.New("http not launched")
		}
		if !httpsLaunched {
			return errors.New("https not launched")
		}
		if email != "test@app.com" {
			return errors.New("email not set")
		}
		if len(domains) != 3 {
			return errors.New("domains not set")
		}
		if len(apps) != 3 {
			return errors.New("apps not set")
		}
		return nil
	}

	return checkData
}

func SubTestHttpsLaunched(t *testing.T, broker *message.MessageBroker) {
	t.Log("SubTestHttpsLaunched")

	check := AddBrokerSubscriberIfHttps(t, broker)
	//  Create  new service
	srv := services.New(broker, nil)

	// Start platform
	err := srv.StartPlatform()
	if err != nil {
		t.Error(err)
	}

	// Check
	err = check()
	if err != nil {
		t.Error(err)
	}
}
