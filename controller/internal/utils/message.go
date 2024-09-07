package utils

import (
	"controller/internal/entities"
	"encoding/json"
	"errors"
)

func ParseMessage(body []byte) error {
	data := new(entities.ResponseMessage)
	err := json.Unmarshal(body, data)
	if err != nil {
		return err
	}

	if data.Status != "success" {
		return errors.New(data.Message)
	}
	return nil
}
