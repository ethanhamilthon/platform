package api

import (
	"encoding/json"
	"net/http"
)

func (a *ApiServer) CheckDomain(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"works": true})
}

func (a *ApiServer) LoginHandler(w http.ResponseWriter, r *http.Request) {
	token, err := a.auth.Login(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (a *ApiServer) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	token, err := a.auth.Register(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (a *ApiServer) MeHandler(w http.ResponseWriter, r *http.Request) {
	username, err := a.auth.Me(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"username": username})
}
