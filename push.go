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
}

func New(clientid, clientSecret string, isdebug bool) *HwPush {
	modeDebug = isdebug

	return &HwPush{
		GrantType:    GrantType,
		ClientId:     clientid,
		ClientSecret: clientSecret,
		client:       &http.Client{},
	}
}

func (this *HwPush) send(req *request) (*response, error) {
	r, err := http.NewRequest(http.MethodPost, req.uri, req.body)
	if err != nil {
		return nil, err
	}

	r.Header.Set("User-Agent", req.ua)
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
	result, err := newResponse(resp.Body)
	return result, err
}

func (this *HwPush) Single(deviceToken, title, content string, custom map[string]string) error {
	err := this.Auth()
	if err != nil {
		return err
	}

	req, err := newRequest(title, content, this.token.Value, []string{deviceToken}, custom, this.ClientId)
	if err != nil {
		return err
	}
	resp, err := this.send(req)
	if err != nil {
		return err
	}
	if resp.Code != 0 {
		return errors.New(fmt.Sprintf("[%d]%s", resp.Code, resp.Message))
	}

	return nil
}

func (this *HwPush) Group(tokens []string, title, content string, custom map[string]string) error {
	err := this.Auth()
	if err != nil {
		return err
	}

	req, err := newRequest(title, content, this.token.Value, tokens, custom, this.ClientId)
	if err != nil {
		return err
	}
	resp, err := this.send(req)
	if err != nil {
		return err
	}
	if resp.Code != 0 || resp.Error != "" {
		return errors.New(fmt.Sprintf("[%d]%s %s", resp.Code, resp.Message, resp.Error))
	}

	return nil
}
