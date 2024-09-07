package tests

import (
	"encoding/json"
	"testing"
)

func DomainTest(t *testing.T) {
	// List has to be empty
	domainListTest(t, 0)

	// Add domain

	domainAddTest(t, "example1.com", false)
	domainAddTest(t, "example2.com", false)
	domainAddTest(t, "example3.com", false)
	domainAddTest(t, "example4.com", false)
	domainAddTest(t, "example4.com", true)
}

func domainAddTest(t *testing.T, domain string, shouldError bool) {
	// Add domain
	body, err := doHttpRequest("http://localhost:8000/api/domain/add", "POST",
		[]byte(`{"domain": "`+domain+`"}`))

	if shouldError {
		if err == nil {
			t.Errorf("error adding domain: %v", err)
		}
	} else {
		if err != nil {
			t.Errorf("error adding domain: %v", err)
		}
	}

	if len(body) != 0 {
		t.Errorf("error adding domain: not empty body")
	}
}

func domainListTest(t *testing.T, wantLen int) {
	body, err := doHttpRequest("http://localhost:8000/api/domain/list", "GET", []byte(`{"token": "`+token+`"}`))

	if err != nil {
		t.Errorf("error listing domains: %v", err)
	}

	if len(body) == 0 {
		t.Errorf("error listing domains: empty body")
	}

	data := make([]struct {
		Domain string `json:"domain"`
	}, 0)

	err = json.Unmarshal(body, &data)
	if err != nil {
		t.Errorf("error listing domains: %v, body: %v", err, string(body))
	}

	if len(data) != wantLen {
		t.Errorf("error listing domains: empty data")
	}
}
