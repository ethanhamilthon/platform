package entities

type ResponseMessage struct {
  Status string `json:"status"` // success, error
  Message string `json:"message"`
}
