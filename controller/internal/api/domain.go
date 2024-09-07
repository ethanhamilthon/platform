package api

import (
	"encoding/json"
	"net/http"
)

func (a *ApiServer) DomainListHandler(w http.ResponseWriter, r *http.Request) {
  domains, err := a.domain.List()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(domains)
}

type DomainAddBody struct {
  Domain string `json:"domain"`
}
func (a *ApiServer) DomainAddHandler(w http.ResponseWriter, r *http.Request) {
  body := new(DomainAddBody)
  err := json.NewDecoder(r.Body).Decode(body)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  err = a.domain.Add(body.Domain)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.WriteHeader(http.StatusOK)
}
