package api

import (
	"controller/internal/services/application"
	"encoding/json"
	"net/http"
)

func (a *ApiServer) DetectRepoHandler(w http.ResponseWriter, r *http.Request) {
	// Parse body
	body := new(application.DetectRepoOptions)
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Call service
	data, err := a.app.DetectRepo(*body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
