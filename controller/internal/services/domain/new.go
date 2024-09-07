package domain

import (
	"controller/internal/db"
	"encoding/json"
	"errors"
)

type DomainService struct {
	db *db.DB
}

func New(db *db.DB) *DomainService {
	return &DomainService{db: db}
}

type DomainData struct {
	Domain string `json:"domain"`
}

func (d *DomainService) Add(domain string) error {
	body, err := d.db.Get("domains")
	if err != nil {
		return err
	}

	data := make([]DomainData, 0)
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	for _, v := range data {
		if v.Domain == domain {
			return errors.New("domain already exists")
		}
	}
	data = append(data, DomainData{Domain: domain})
	body, err = json.Marshal(data)
	if err != nil {
		return err
	}
	return d.db.Set("domains", body)
}

func (d *DomainService) List() ([]DomainData, error) {
	body, err := d.db.Get("domains")
	if err != nil {
		return nil, err
	}

	data := make([]DomainData, 0)
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
