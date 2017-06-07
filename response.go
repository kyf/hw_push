package hw_push

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
)

const (
	RESP_STATE_OK    = 1
	RESP_STATE_ERROR = 0
)

type authResponse struct {
	Error       int    `json:"error,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
	ExpireIn    int    `json:"expire_in,omitempty"`
	ErrorDes    string `json:"error_description,omitempty"`
}

type response struct {
	Code    int    `json:"resultcode,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	ReqId   string `json:"requestID,omitempty"`
}

func newResponse(reader io.Reader) (*response, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	if modeDebug {
		log.Printf("response : %s", string(data))
	}
	resp := response{}
	decoder := json.NewDecoder(bytes.NewReader(data))
	err = decoder.Decode(&resp)
	return &resp, err
}
