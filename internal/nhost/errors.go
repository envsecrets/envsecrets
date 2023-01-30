package nhost

import "encoding/json"

type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Code    string `json:"Message"`
}

func New(data []byte) *Error {
	var response Error
	json.Unmarshal(data, &response)
	return &response
}
