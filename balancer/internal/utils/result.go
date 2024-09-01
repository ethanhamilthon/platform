package utils

import "encoding/json"

type Message struct {
	M string `json:"message"`
}

func Success() ([]byte, error) {
	result, err := json.Marshal(Message{
		M: "success",
	})
	if err != nil {
		return []byte{}, err
	}

	return result, nil
}

func Error(err error) []byte {
	result, _ := json.Marshal(Message{
		M: err.Error(),
	})

	return result
}
