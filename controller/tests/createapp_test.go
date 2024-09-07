package tests

import (
	"controller/internal/config"
	"controller/internal/entities"
	"controller/internal/message"
	"controller/internal/services/application"
	"controller/internal/utils"
	"encoding/json"
	"testing"

	"github.com/nats-io/nats.go/jetstream"
)

func createAppTest(t *testing.T) {
	if token == "" {
		t.Fatalf("error create app: empty token")
	}
  stop := mockStream(t)
  defer stop()
	createAppRequest(t)
}

func createAppRequest(t *testing.T) {
	// pull repo from github
	pulldata := application.PublicGithubPullOptions{
		Url:    "https://github.com/ethanmotion/getstarted",
		Branch: "main",
		Root:   "/",
	}
	body, err := doHttpRequest("http://localhost:8000/api/github/pull", "POST", pulldata)
	if err != nil {
		t.Fatalf("error create app: %v, data: %v", err, string(body))
	}

	data := make(map[string]string)
	err = json.Unmarshal(body, &data)
	if err != nil {
		t.Fatalf("error create app: %v, data: %v", err, string(body))
	}
	if data["repo_id"] == "" {
		t.Fatalf("error create app: %v", err)
	}

	// detect repo
	body, err = doHttpRequest("http://localhost:8000/api/repo/detect", "POST", map[string]string{"repo_id": data["repo_id"]})
	if err != nil {
		t.Fatalf("error create app: %v, data: %v", err, string(body))
	}
	detect_data := application.DetectRepoResponse{}
	err = json.Unmarshal(body, &detect_data)
	if err != nil {
		t.Fatalf("error create app: %v", err)
	}

	// create app
	app_data := application.CreateAppOptions{
		Name:        "testapp",
		Description: "test app",
		Source: application.SourceData{
			From: "github",
			URL:  pulldata.Url,
			Data: "",
		},
		Network: application.NetworkData{
			Hostname:   "app.com",
			PathPrefix: "/api",
			Port:       detect_data.Port,
		},
		Build: application.BuildData{
			DockerfileExits: detect_data.DockerfileExits,
			Dockerfile:      "",
			Language:        detect_data.Language,
			Framework:       detect_data.Framework,
			InstallCmds:     detect_data.InstallCmds,
			BuildCmds:       detect_data.BuildCmds,
			RunCmds:         detect_data.RunCmds,
			Dependencies:    detect_data.Dependencies,
		},
		Envs:   detect_data.Envs,
		Binds:  []string{},
		RepoID: data["repo_id"],
	}

	body, err = doHttpRequest("http://localhost:8000/api/app/create", "POST", app_data)
	if err != nil {
		t.Fatalf("error create app: %v, data: %v", err, string(body))
	}
	data = make(map[string]string)
	err = json.Unmarshal(body, &data)
	if err != nil {
		t.Fatalf("error create app: %v, data: %v", err, string(body))
	}
	if data["app_id"] == "" {
		t.Fatalf("error create app: %v", err)
	}

}

func mockStream(t *testing.T) func() {
	// Create stream
	broker, err := message.New()
	if err != nil {
		t.Fatalf("error init nats: %v", err)
	}
  stream, err := broker.CreateStream(message.CreateStramOptions{
		Name: config.AdapterStream,
		Subjects: []string{
			config.AdapterStreamPathPrefix,
		},
	})
	if err != nil {
		t.Fatalf("error init nats: %v", err)
	}

	// Consume to build image
	bi_consumer, err := broker.CreateConsumer(message.CreateConsumerOptions{
		Name:    "adapter:buildimage",
		Subject: config.AdapterBuildImage,
	},stream, func(m jetstream.Msg) {
		t.Log("consume build image")
		res_data := new(application.BuildImageOptions)
		if err := json.Unmarshal(m.Data(), res_data); err != nil {
			t.Fatalf("error consume build image: %v", err)
		}
		if err := m.Ack(); err != nil {
			t.Fatalf("error consume build image: %v", err)
		}
	})

	// Consume to pull github
	pi_consumer, err := broker.CreateConsumer(message.CreateConsumerOptions{
		Name:    "adapter-pullgithubs",
		Subject: config.AdapterPullGithub,
	},stream, func(m jetstream.Msg) {
		t.Log("consume pull github")
		id := m.Headers().Get("Nats-Msg-Id")
		if id == "" {
			t.Fatalf("error consume pull image: empty id")
		}
		data := make(map[string]string)
		if err := json.Unmarshal(m.Data(), &data); err != nil {
			t.Fatalf("error consume pull image: %v", err)
		}
		if err := m.Ack(); err != nil {
			t.Fatalf("error consume pull image: %v", err)
		}

		// Send error response
		if data["url"] == "" || data["branch"] == "" || data["root"] == "" {
			if err = broker.PublishID(config.ControllerAnswer, entities.ResponseMessage{
				Message: "empty fields",
				Status:  "error",
			}, id); err != nil {
				t.Fatalf("error consume pull image: %v", err)
			}
		}

		// Send success response
		if err = broker.PublishID(config.ControllerAnswer, entities.ResponseMessage{
			Message: "success",
			Status:  "success",
		}, id); err != nil {
			t.Fatalf("error consume pull image: %v", err)
		}
	})

	// Consume to detect
	d_consumer, err := broker.CreateConsumer(message.CreateConsumerOptions{
		Name:    "adapter:detect",
		Subject: config.AdapterDetectRepo,
	},stream, func(m jetstream.Msg) {
		t.Log("consume detect image")
		id := m.Headers().Get("Nats-Msg-Id")
		if id == "" {
			t.Fatalf("error consume detect image: empty id")
		}

		// get data
		request_data := make(map[string]string)
		if err := json.Unmarshal(m.Data(), &request_data); err != nil {
			t.Fatalf("error consume detect image: %v", err)
		}
		if request_data["repo_id"] == "" {
			t.Fatalf("error consume detect image: empty or not matched repo_id")
		}
		if err := m.Ack(); err != nil {
			t.Fatalf("error consume detect image: %v", err)
		}

		// prepare response
		res_data := application.DetectRepoResponse{
			DockerfileExits: true,
			Language:        "",
			Framework:       "",
			InstallCmds:     []string{},
			BuildCmds:       []string{},
			RunCmds:         []string{},
			Dependencies:    []string{},
			Envs:            []string{"ENV1=VALUE1"},
			Port:            "3000",
		}
		if isThereDockerfile := utils.FlipCoin(); !isThereDockerfile {
			res_data.DockerfileExits = false
			res_data.Language = "node:20"
			res_data.Framework = "express"
			res_data.InstallCmds = []string{"npm install"}
			res_data.BuildCmds = []string{"npm run build"}
			res_data.RunCmds = []string{"npm run start"}
			res_data.Dependencies = []string{}
			res_data.Envs = []string{"ENV1=VALUE1"}
			res_data.Port = "3000"
		}

		// send response
		if err := broker.PublishID(config.ControllerAnswer, res_data, id); err != nil {
			t.Fatalf("error consume detect image: %v", err)
		}

	})

	// Consume to create dockerfile
	cd_consumer, err := broker.CreateConsumer(message.CreateConsumerOptions{
		Name:    "adapter:create-dockerfile",
		Subject: config.AdapterCreateDockerfile,
	},stream, func(m jetstream.Msg) {
		t.Log("consume create dockerfile")
		id := m.Headers().Get("Nats-Msg-Id")
		if id == "" {
			t.Fatalf("error consume create dockerfile: empty id")
		}

		// get data
		request_data := new(application.CreateDockerfileOptions)
		if err := json.Unmarshal(m.Data(), request_data); err != nil {
			t.Fatalf("error consume create dockerfile: %v", err)
		}
		if err := m.Ack(); err != nil {
			t.Fatalf("error consume create dockerfile: %v", err)
		}
		if request_data.Language != "node:20" {
			t.Fatalf("error consume create dockerfile: empty language")
		}
		// send response
		if err := broker.PublishID(config.ControllerAnswer, map[string]string{
			"status":     "success",
			"dockerfile": "some dockerfile",
		}, id); err != nil {
			t.Fatalf("error consume create dockerfile: %v", err)
		}
	})
  stop := func() {
    bi_consumer.Stop()
    d_consumer.Stop()
    cd_consumer.Stop()
    pi_consumer.Stop()
  }
  return stop
}
