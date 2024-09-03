package tests

import (
	"encoding/json"
	"errors"
	"manager/internal/message"
	"manager/internal/services"
	"testing"
)

func AddBrokerSubscriber(t *testing.T, broker *message.MessageBroker) func() error {
	apps := make([]services.Application, 0)
	httpLaunched := false

	broker.Response("controller:get:balancer-state", func(data []byte) ([]byte, error) {
		data, err := json.Marshal(map[string]interface{}{
			"domain": []string{},
			"email":  "",
		})
		if err != nil {
			t.Error(err)
			return []byte{}, err
		}
		return data, nil
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

	broker.Response("balancer:launch:http", func(data []byte) ([]byte, error) {
		httpLaunched = true
		t.Log("http listening")
		return SuccessMessage(), nil
	})

	checkData := func() error {
		controllerRunned := false
		controllerApp := services.Application{
			Domain:        "*",
			ContainerName: "controller",
			Port:          "8000",
			PathPrefix:    "/api",
		}

		uiRunned := false
		uiApp := services.Application{
			Domain:        "*",
			ContainerName: "aranea-ui",
			Port:          "3000",
			PathPrefix:    "/",
		}

		if !httpLaunched {
			return errors.New("http not launched")
		}

		if len(apps) != 2 {
			return errors.New("controller and ui not added")
		}

		for _, app := range apps {
			if app.ContainerName == controllerApp.ContainerName {
				if app.Domain != controllerApp.Domain || app.Port != controllerApp.Port || app.PathPrefix != controllerApp.PathPrefix {
					return errors.New("controller not added")
				}

				controllerRunned = true
			}

			if app.ContainerName == uiApp.ContainerName {
				if app.Domain != uiApp.Domain || app.Port != uiApp.Port || app.PathPrefix != uiApp.PathPrefix {
					return errors.New("ui not added")
				}

				uiRunned = true
			}
		}

		if !controllerRunned {
			return errors.New("controller not added")
		}

		if !uiRunned {
			return errors.New("ui not added")
		}

		return nil
	}

	return checkData
}

func SubTestHttpLaunched(t *testing.T, broker *message.MessageBroker) {
	t.Log("SubTestOnlyHttpLaunched")

	check := AddBrokerSubscriber(t, broker)
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
