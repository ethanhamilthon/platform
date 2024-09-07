package application

import (
	"controller/internal/db"
	"controller/internal/message"
	"controller/internal/services/listener"
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type ApplicationService struct {
	db       *db.DB
	broker   *message.Broker
	listener *listener.ListenerService
}

func New(db *db.DB, broker *message.Broker, listener *listener.ListenerService) *ApplicationService {
	return &ApplicationService{
		db:       db,
		broker:   broker,
		listener: listener,
	}
}

type AppData struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Source      SourceData    `json:"source"`
	State       string        `json:"state"`
	Build       BuildData     `json:"build"`
	Image       ImageData     `json:"image"`
	Container   ContainerData `json:"container"`
	Network     NetworkData   `json:"network"`
	RepoID      string        `json:"repo_id"`
}

type ImageData struct {
	ImageName         string   `json:"image_name"`
	ImageVersions     []string `json:"image_versions"`
	LastStableVersion string   `json:"last_stable_version"`
}

type ContainerData struct {
	ContainerTag      string   `json:"container_name"`
	ContainerVersions []string `json:"container_versions"`
	LastStableVersion string   `json:"last_stable_version"`
	Envs              []string `json:"envs"`
	Binds             []string `json:"binds"`
}

type SourceData struct {
	From string `json:"from" validate:"required"`
	URL  string `json:"url" validate:"required"`
	Data string `json:"data"`
}

type NetworkData struct {
	Hostname   string `json:"hostname"`
	PathPrefix string `json:"path_prefix"`
	Port       string `json:"port" validate:"required"`
}

type BuildData struct {
	DockerfileExits bool     `json:"dockerfile_exits"`
	Dockerfile      string   `json:"dockerfile"`
	Language        string   `json:"language"`
	Framework       string   `json:"framework"`
	InstallCmds     []string `json:"install_cmds"`
	BuildCmds       []string `json:"build_cmds"`
	RunCmds         []string `json:"run_cmds"`
	Dependencies    []string `json:"dependencies"`
}

type CreateAppOptions struct {
	Name        string      `json:"name" validate:"required"`
	Description string      `json:"description" validate:"required"`
	Source      SourceData  `json:"source"`
	Network     NetworkData `json:"network"`
	Build       BuildData   `json:"build"`
	Envs        []string    `json:"envs"`
	Binds       []string    `json:"binds"`
	RepoID      string      `json:"repo_id" validate:"required"`
}

func (s *ApplicationService) Create(opts CreateAppOptions) (string, error) {
	appID := uuid.New().String()
	app := AppData{
		ID:          appID,
		Name:        opts.Name,
		Description: opts.Description,
		Source:      opts.Source,
		State:       "creating",
		Build:       opts.Build,
		Image: ImageData{
			ImageName:         strings.ReplaceAll(strings.ToLower(opts.Name), " ", "-") + "-image",
			ImageVersions:     []string{"1"},
			LastStableVersion: "",
		},
		Container: ContainerData{
			ContainerTag:      strings.ReplaceAll(strings.ToLower(opts.Name), " ", "-") + "-container",
			ContainerVersions: []string{"1"},
			LastStableVersion: "",
			Envs:              opts.Envs,
			Binds:             opts.Binds,
		},
		Network: opts.Network,
		RepoID:  opts.RepoID,
	}
	// switch type source
	switch app.Source.From {
	case "github":
		return s.createPublicGithubApp(app)
	case "gitlab":
		return "", errors.New("gitlab not implemented")
	case "dockerhub":
		return "", errors.New("dockerhub not implemented")
	case "github_private":
		return "", errors.New("github private not implemented")
	case "gitlab_private":
		return "", errors.New("gitlab private not implemented")
	case "dockerhub_private":
		return "", errors.New("dockerhub private not implemented")
	default:
		return "", errors.New("invalid source")
	}
}

func (s *ApplicationService) saveApp(app AppData) error {
	// Save app
	appBody, err := json.Marshal(app)
	if err != nil {
		return err
	}
	err = s.db.Set("app"+app.ID, appBody)
	if err != nil {
		return err
	}
	return nil
}

type BuildImageOptions struct {
	AppID           string `json:"app_id"`
	DockerfileExits bool   `json:"dockerfile_exits"`
	Dockerfile      string `json:"dockerfile"`
	ImageName       string `json:"image_name"`
	RepoID          string `json:"repo_id"`
}

type CreateDockerfileOptions struct {
	Language     string   `json:"language"`
	Framework    string   `json:"framework"`
	InstallCmds  []string `json:"install_cmds"`
	BuildCmds    []string `json:"build_cmds"`
	RunCmds      []string `json:"run_cmds"`
	Dependencies []string `json:"dependencies"`
}

type DockerfileResponse struct {
	Dockerfile string `json:"dockerfile"`
	Status     string `json:"status"`
}
