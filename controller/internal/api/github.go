package api

import (
	"controller/internal/services/application"
	"encoding/json"
	"net/http"
)

func (a *ApiServer) PublicGithubPullHandler(w http.ResponseWriter, r *http.Request) {
	// Parse body
	body := new(application.PublicGithubPullOptions)
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Call service
	repoID, err := a.app.PullGithubPublic(*body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"repo_id": repoID})
}
