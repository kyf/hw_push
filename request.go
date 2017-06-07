package hw_push

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"
	"time"
)

const (
	UserAgent = "liurenyou_push_client_v0.0.1"
)

type request struct {
	body io.Reader
	uri  string
	ua   string
}

func newRequest(params url.Values, token, method string) (*request, error) {
	data := url.Values{
		"access_token": []string{token},
		"nsp_fmt":      []string{ResponseFormat},
		"nsp_ts":       []string{fmt.Sprintf("%d", time.Now().Unix())},
		"nsp_svc":      []string{method},
		//"userType":     []string{"-1"},
		"device_type": []string{"android"},
		"expire_time": []string{time.Now().Add(time.Second * 60 * 60 * 24 * 2).Format("2006-01-02 15:04")},
	}

	for k, _ := range params {
		data.Set(k, params.Get(k))
	}

	if modeDebug {
		log.Printf("request : %+v", data)
	}

	return &request{body: strings.NewReader(data.Encode()), uri: PushUri, ua: UserAgent}, nil
}
