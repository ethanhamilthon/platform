package builder

import (
	"controller/internal/config"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
)

type SourceOptions struct {
	From   string `json:"from"`
	Url    string `json:"url"`
	Branch string `json:"branch"`
	Root   string `json:"root"`
}

type NetworkOptions struct {
	Hostname   string `json:"hostname"`
	PathPrefix string `json:"path_prefix"`
	Port       int    `json:"port"`
}

type BuildOptions struct {
	DockerfileExists string   `json:"dockerfile_exists"`
	Language         string   `json:"language"`
	Framework        string   `json:"framework"`
	InstallCmds      []string `json:"install_cmds"`
	BuildCmds        []string `json:"build_cmds"`
	RunCmds          []string `json:"run_cmds"`
}

type RunOptions struct {
	Envs  []string `json:"envs"`
	Binds []string `json:"binds"`
}

type WithCodeOptions struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Source      SourceOptions  `json:"source"`
	Network     NetworkOptions `json:"network"`
	Build       BuildOptions   `json:"build"`
	Run         RunOptions     `json:"run"`
}

func (b *BuilderService) NewAppWithCode(m jetstream.Msg) {
	data := WithCodeOptions{}
	err := json.Unmarshal(m.Data(), &data)
	if err != nil {
		return
	}
	m.Ack()

	// Clone repo
	repoID := uuid.New().String()
	msgID := uuid.New().String()
	wait := b.listener.Add(msgID)
	err = b.broker.PublishID(config.AdapterPullGithub, map[string]any{
		"id":          repoID,
		"url":         data.Source.Url,
		"branch":      data.Source.Branch,
		"root":        data.Source.Root,
		"temporarily": false,
	}, msgID)
	if err != nil {
		return
	}
	body, err := wait()
	if err != nil {
		return
	}
	// Detect repo

	// Add watcher

	// Create app

	// Run app

	// If Runs add to stable and add to balancer

}
