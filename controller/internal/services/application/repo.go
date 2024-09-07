package application

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

type DetectRepoOptions struct {
	RepoId string `json:"repo_id"`
}

type DetectRepoResponse struct {
	DockerfileExits bool     `json:"dockerfile_exits"`
	Language        string   `json:"language"`
	Framework       string   `json:"framework"`
	InstallCmds     []string `json:"install_cmds"`
	BuildCmds       []string `json:"build_cmds"`
	RunCmds         []string `json:"run_cmds"`
	Dependencies    []string `json:"dependencies"`
	Envs            []string `json:"envs"`
	Port            string   `json:"port"`
}

func (s *ApplicationService) DetectRepo(opts DetectRepoOptions) (*DetectRepoResponse, error) {
	msgid := uuid.New().String()
	wait := s.listener.Add(msgid)
	if err := s.broker.PublishID("adapter.detect.repo", opts, msgid); err != nil {
		return nil, errors.New("publish failed")
	}
	body, err := wait()
	if err != nil {
		return nil, errors.New("detect failed")
	}

	data := new(DetectRepoResponse)
	if err := json.Unmarshal(body, data); err != nil {
		return nil, errors.New("get data failed")
	}

	return data, nil
}
