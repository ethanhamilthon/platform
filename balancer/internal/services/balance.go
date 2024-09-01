package services

import (
	"balancer/internal/utils"
	"encoding/json"
	"errors"
	"strings"
)

type AddAppBody struct {
	Domain    string `json:"domain"`
	Container string `json:"container"`
	Port      string `json:"port"`
	Path      string `json:"path"`
}

func (s *Service) AddApp(data []byte) ([]byte, error) {
	// Json body
	var body AddAppBody
	err := json.Unmarshal(data, &body)
	if err != nil {
		return []byte{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// If there is already domain in system
	for i, dom := range s.domains {
		if dom.Domain == body.Domain {
			s.domains[i].Apps = append(s.domains[i].Apps, AppType{
				Port:          body.Port,
				Path:          body.Path,
				ContainerName: body.Container,
			})
			return utils.Success()
		}
	}

	// if there isnt domain, we will add new domain
	s.domains = append(s.domains, DomainType{
		Domain: body.Domain,
		Apps: []AppType{
			{
				Port:          body.Port,
				Path:          body.Path,
				ContainerName: body.Container,
			},
		},
	})

	return utils.Success()
}

type ChangeContainerBody struct {
	Original string `json:"original_name"`
	New      string `json:"new_name"`
}

func (s *Service) ChangeContainer(data []byte) ([]byte, error) {
	// Body
	var body ChangeContainerBody
	err := json.Unmarshal(data, &body)
	if err != nil {
		return []byte{}, err
	}

	// Change container name
	for i, v := range s.domains {
		for j, a := range v.Apps {
			if body.Original == a.ContainerName {
				s.domains[i].Apps[j].ContainerName = body.New
			}
		}
	}

	return utils.Success()
}

type RemoveAppBody struct {
	Domain    string `json:"domain"`
	Container string `json:"container"`
}

func (s *Service) RemoveApp(data []byte) ([]byte, error) {
	// body
	var body RemoveAppBody
	err := json.Unmarshal(data, &body)
	if err != nil {
		return []byte{}, err
	}

	// Create new array with only not deleted apps
	apps := make([]AppType, 0)
	for i, v := range s.domains {
		if body.Domain == v.Domain {
			for _, a := range v.Apps {
				if body.Container != a.ContainerName {
					apps = append(apps, a)
				}
			}

			s.domains[i].Apps = apps
		}

	}
	return utils.Success()
}

// Returns an url to service in docker network by domain name and path
func (s *Service) getServiceUrl(host string, path string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var matchedApp *AppType
	var longestMatch int
	for _, dom := range s.domains {
		if dom.Domain == host {
			for _, app := range dom.Apps {
				if strings.HasPrefix(path, app.Path) {
					if len(app.Path) > longestMatch {
						longestMatch = len(app.Path)
						matchedApp = &app
					}
				}
			}
		}
	}

	if matchedApp == nil {
		return "", errors.New("not found")
	}
	return "http://" + matchedApp.ContainerName + ":" + matchedApp.Port, nil
}
