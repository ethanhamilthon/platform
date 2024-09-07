package api

import (
	"controller/internal/services/application"
	"encoding/json"
	"net/http"
)

func (a *ApiServer) CreateApp(w http.ResponseWriter, r *http.Request) {
	body := new(application.CreateAppOptions)
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	app_id, err := a.app.Create(*body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"app_id": app_id})
}
