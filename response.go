package hw_push

import (
	"io"
)

const (
	RESP_STATE_OK    = 1
	RESP_STATE_ERROR = 0
)

type response struct {
	Code    int    `json:"resultcode,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	ReqId   string `json:"requestID,omitempty"`
}

func newResponse(reader io.Reader) *response {
	return &response{}
}
