package hw_push

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func newRequestV1(title, content, token string, devices []string, custom map[string]string, appid string) (*request, error) {

	intent := fmt.Sprintf("intent://com.liurenyou.im/notification?action=%s#Intent;scheme=wonderfullpush;launchFlags=0x10000000;end", custom["orderid"])

	payload := map[string]interface{}{
		"validate_only": false,
		"message": map[string]interface{}{
			"notification": map[string]interface{}{
				"title": title,
				"body":  content,
			},
			"android": map[string]interface{}{
				"notification": map[string]interface{}{
					"title": title,
					"body":  content,
					"click_action": map[string]interface{}{
						"type":   1,
						"intent": intent,
					},
				},
			},
			"token": devices,
		},
	}

	_payload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf(PushUriV1, appid)
	return &request{body: bytes.NewReader(_payload), uri: uri, ua: UserAgent}, nil
}

func newRequest(title, content, token string, devices []string, custom map[string]string, appid string) (*request, error) {
	data := url.Values{
		"access_token": []string{token},
		"nsp_ts":       []string{fmt.Sprintf("%d", time.Now().Unix())},
		"nsp_svc":      []string{MethodName},
		"expire_time":  []string{time.Now().Add(time.Second * 60 * 60 * 24 * 2).Format("2006-01-02T15:04")},
	}

	custom["text"] = content
	custom["title"] = title

	payload := map[string]interface{}{
		"hps": map[string]interface{}{
			"msg": map[string]interface{}{
				"type": 1,
				"body": custom,
			},
		},
	}

	_payload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	data.Add("device_token_list", "[\""+strings.Join(devices, "\", \"")+"\"]")
	data.Add("payload", string(_payload))

	query := make(url.Values)
	query.Add("nsp_ctx", `{"ver":"1", "appId":"`+appid+`"}`)

	uri := PushUri + "?" + query.Encode()
	return &request{body: strings.NewReader(data.Encode()), uri: uri, ua: UserAgent}, nil
}
