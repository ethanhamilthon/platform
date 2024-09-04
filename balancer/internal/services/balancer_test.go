package services

import (
	"encoding/json"
	"testing"
)

func testapps() (*Service, error) {
	s := New()
	apps := []AddAppBody{
		{
			Domain:    "*",
			Container: "controller",
			Port:      "8000",
			Path:      "/api",
		}, {
			Domain:    "*",
			Container: "aranea-ui",
			Port:      "3000",
			Path:      "/",
		}, {
			Domain:    "app.com",
			Container: "blogapp",
			Port:      "4000",
			Path:      "/",
		}, {
			Domain:    "app.com",
			Container: "blogapp-api",
			Port:      "4000",
			Path:      "/api",
		},
		{
			Domain:    "*.some.com",
			Container: "some-ui",
			Port:      "5000",
			Path:      "/",
		}, {
			Domain:    "*.some.com",
			Container: "some-api",
			Port:      "5002",
			Path:      "/api",
		},
	}

	for _, app := range apps {
		data, err := json.Marshal(app)
		if err != nil {
			return s, err
		}
		s.AddApp(data)
	}
	return s, nil
}

func TestNoApps(t *testing.T) {
	s := New()
	url, err := s.getServiceUrl("app.com", "/api/users")
	if err == nil {
		t.Errorf("expected error, got: %v", url)
		return
	}
	if err.Error() != "not found" {
		t.Errorf("expected not found, got %s", err.Error())
	}
	_, err = s.getServiceUrl("localhost:8080", "/")
	if err == nil {
		t.Error("expected error")
		return
	}
	if err.Error() != "not found" {
		t.Errorf("expected not found, got %s", err.Error())
	}
}

func TestBalancerDomains(t *testing.T) {
	s, err := testapps()
	if err != nil {
		t.Error(err)
	}
	if len(s.domains) != 3 {
		t.Error("expected 3 domains")
	}

	url, err := s.getServiceUrl("app.com", "/api/users")
	if err != nil {
		t.Error(err)
	}
	if url != "http://blogapp-api:4000" {
		t.Error("expected http://blogapp-api:4000/api got:" + url)
	}

	url, err = s.getServiceUrl("app.com", "/about")
	if err != nil {
		t.Error(err)
	}
	if url != "http://blogapp:4000" {
		t.Error("expected http://blogapp:4000/ got:" + url)
	}

	url, err = s.getServiceUrl("another.com", "/mypage/ao")
	if err != nil {
		t.Error(err)
	}
	if url != "http://aranea-ui:3000" {
		t.Error("expected http://aranea-ui:3000/ got:" + url)
	}

	url, err = s.getServiceUrl("another.com", "/api/docker")
	if err != nil {
		t.Error(err)
	}
	if url != "http://controller:8000" {
		t.Error("expected http://controller:8000/ got:" + url)
	}

	url, err = s.getServiceUrl("api.some.com", "/dashboard?i=doesnt%20matter")
	if err != nil {
		t.Error(err)
	}
	if url != "http://some-ui:5000" {
		t.Error("expected http://some-ui:5000/ got:" + url)
	}

	url, err = s.getServiceUrl("api.some.com", "/api/some/thing?i=doesnt%20matter")
	if err != nil {
		t.Error(err)
	}
	if url != "http://some-api:5002" {
		t.Error("expected http://some-api:5002/ got:" + url)
	}

	url, err = s.getServiceUrl("some.com", "/")
	if err != nil {
		t.Error(err)
	}
	if url != "http://aranea-ui:3000" {
		t.Error("expected http://aranea-ui:3000/ got:" + url)
	}
}
