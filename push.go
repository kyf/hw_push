package hw_push

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"sync"
	"time"
)

const (
	ExpireBuff = 60 * 30
	DebugOn    = true
	DebugOff   = false
)

var (
	modeDebug = false
)

type AccessToken struct {
	Value    string
	ExpireIn time.Time
}

type HwPush struct {
	GrantType    string
	ClientId     string
	ClientSecret string
	client       *http.Client
	token        *AccessToken
	tokenLocker  sync.Mutex
	Version      string
}

func New(clientid, clientSecret string, isdebug bool, version string) *HwPush {
	modeDebug = isdebug

	return &HwPush{
		GrantType:    GrantType,
		ClientId:     clientid,
		ClientSecret: clientSecret,
		client:       &http.Client{},
		Version:      version,
	}
}

func (this *HwPush) send(req *request) (interface{}, error) {
	r, err := http.NewRequest(http.MethodPost, req.uri, req.body)
	if err != nil {
		return nil, err
	}

	r.Header.Set("User-Agent", req.ua)
	r.Header.Set("Authorization", this.token.Value)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if modeDebug {
		reqBytes, _ := httputil.DumpRequest(r, true)
		log.Printf("request : %s", string(reqBytes))
	}

	resp, err := this.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return newResponseV1(resp.Body)
	/*
		if this.Version == "v1" {
			return newResponseV1(resp.Body)
		} else {
			return newResponse(resp.Body)
		}
	*/
}

func (this *HwPush) Single(deviceToken, title, content string, custom map[string]string) error {
	err := this.Auth()
	if err != nil {
		return err
	}

	var req *request
	if this.Version == "v1" {
		req, err = newRequestV1(title, content, this.token.Value, []string{deviceToken}, custom, this.ClientId, "")
	} else {
		req, err = newRequest(title, content, this.token.Value, []string{deviceToken}, custom, this.ClientId)
	}
	if err != nil {
		return err
	}
	resp, err := this.send(req)
	if err != nil {
		return err
	}
	switch _resp := resp.(type) {
	case *responseV1:
		log.Print("v1 response is ", _resp)
		if _resp.Code != V1_SUCCESS_CODE {
			return errors.New(fmt.Sprintf("[%s]%s", _resp.Code, _resp.Message))
		}
	case *response:
		log.Print("v0 response is ", _resp)
		if _resp.Code != 0 {
			return errors.New(fmt.Sprintf("[%d]%s", _resp.Code, _resp.Message))
		}
	}

	return nil
}

func (this *HwPush) Group(tokens []string, title, content string, custom map[string]string) error {
	err := this.Auth()
	if err != nil {
		return err
	}

	req, err := newRequestV1(title, content, this.token.Value, tokens, custom, this.ClientId, "")
	if err != nil {
		return err
	}
	resp, err := this.send(req)
	if err != nil {
		return err
	}
	switch _resp := resp.(type) {
	case *responseV1:
		if _resp.Code != V1_SUCCESS_CODE {
			return errors.New(fmt.Sprintf("[%s]%s ", _resp.Code, _resp.Message))
		}
	case *response:
		if _resp.Code != 0 {
			return errors.New(fmt.Sprintf("[%d]%s\t%s", _resp.Code, _resp.Message, _resp.Error))
		}
	}

	return nil
}

func (this *HwPush) All(title, content string, custom map[string]string) error {
	err := this.Auth()
	if err != nil {
		return err
	}

	req, err := newRequestV1(title, content, this.token.Value, []string{}, custom, this.ClientId, "default")
	if err != nil {
		return err
	}
	resp, err := this.send(req)
	if err != nil {
		return err
	}
	switch _resp := resp.(type) {
	case *responseV1:
		if _resp.Code != V1_SUCCESS_CODE {
			return errors.New(fmt.Sprintf("[%s]%s ", _resp.Code, _resp.Message))
		}
	case *response:
		if _resp.Code != 0 {
			return errors.New(fmt.Sprintf("[%d]%s\t%s", _resp.Code, _resp.Message, _resp.Error))
		}
	}

	return nil
}
