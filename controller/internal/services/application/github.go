package application

import (
	"controller/internal/config"
	"controller/internal/utils"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

type PublicGithubPullOptions struct {
	Url    string `json:"url"`
	Branch string `json:"branch"`
	Root   string `json:"root"`
}

func (s *ApplicationService) PullGithubPublic(opts PublicGithubPullOptions) (string, error) {
	// Send message to pull
	msgid := uuid.New().String()
	wait := s.listener.Add(msgid)
	if err := s.broker.PublishID(config.AdapterPullGithub, opts, msgid); err != nil {
		return "", errors.New("publish failed")
	}
	if body, err := wait(); err != nil || utils.ParseMessage(body) != nil {
		return "", err
	}
	return msgid, nil
}

func (s *ApplicationService) createPublicGithubApp(app AppData) (string, error) {
	// if there is already dockerfile, use it
	if app.Build.DockerfileExits {
		// Send message to build
		if err := s.broker.Publish(config.AdapterBuildImage, BuildImageOptions{
			AppID:           app.ID,
			DockerfileExits: app.Build.DockerfileExits,
			Dockerfile:      app.Build.Dockerfile,
			ImageName:       app.Image.ImageName + ":" + app.Image.ImageVersions[0],
			RepoID:          app.RepoID,
		}); err != nil {
			return "", err
		}
	} else {
		// Create Dockerfile
		msgID := uuid.New().String()
		wait := s.listener.Add(msgID)
		if err := s.broker.PublishID(config.AdapterCreateDockerfile, CreateDockerfileOptions{
			Language:     app.Build.Language,
			Framework:    app.Build.Framework,
			InstallCmds:  app.Build.InstallCmds,
			BuildCmds:    app.Build.BuildCmds,
			RunCmds:      app.Build.RunCmds,
			Dependencies: app.Build.Dependencies,
		}, msgID); err != nil {
			return "", err
		}
		body, err := wait()
		if err != nil {
			return "", err
		}
		data := new(DockerfileResponse)
		if err := json.Unmarshal(body, data); err != nil {
			return "", err
		}
		if data.Status != "success" {
			return "", errors.New("create dockerfile failed")
		}
		app.Build.Dockerfile = data.Dockerfile

		// Send message to build
		if err := s.broker.Publish("adapter.build.image", BuildImageOptions{
			AppID:           app.ID,
			DockerfileExits: app.Build.DockerfileExits,
			Dockerfile:      app.Build.Dockerfile,
			ImageName:       app.Image.ImageName + ":" + app.Image.ImageVersions[0],
			RepoID:          app.RepoID,
		}); err != nil {
			return "", err
		}
	}

	// Save app
	if err := s.saveApp(app); err != nil {
		return "", err
	}
	return app.ID, nil
}
