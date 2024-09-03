package utils

import "encoding/json"

func IsSuccessMessage(data []byte) bool {
	var body struct {
		Message string `json:"message"`
	}

	err := json.Unmarshal(data, &body)
	if err != nil {
		return false
	}

	if body.Message == "success" {
		return true
	}

	return false
}
