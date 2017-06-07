package hw_push

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	ExpireBuff = 60 * 30
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
	resp, err := this.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result, err := newResponse(resp.Body)
	return result, err
}

func (this *HwPush) Single(deviceToken, content string) error {
	err := this.Auth()
	if err != nil {
		return err
	}

	data := url.Values{
		"deviceToken": []string{deviceToken},
		"message":     []string{content},
		"priority":    []string{"0"},
		"cacheMode":   []string{"1"},
		"msgType":     []string{"-1"},
	}

	method := "openpush.message.single_send"

	req, err := newRequest(data, this.token.Value, method)
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

func (this *HwPush) Group(tokens []string, message string) error {
	err := this.Auth()
	if err != nil {
		return err
	}

	data := url.Values{
		"deviceTokenList": []string{strings.Join(tokens, ",")},
		"message":         []string{message},
		"cacheMode":       []string{"1"},
		"msgType":         []string{"-1"},
	}
	method := "openpush.message.batch_send"
	req, err := newRequest(data, this.token.Value, method)
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

func (this *HwPush) All(message string) error {
	err := this.Auth()
	if err != nil {
		return err
	}

	data := url.Values{
		"push_type": []string{"2"},
		"android":   []string{message},
	}
	method := "openpush.openapi.notification_send"
	req, err := newRequest(data, this.token.Value, method)
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
