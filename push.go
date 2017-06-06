package hw_push

import (
	"log"
	"net/http"
)

type HwPush struct {
	GrantType    string
	ClientId     string
	ClientSecret string
	client       *http.Client
}

func New() *HwPush {
	return &HwPush{GrantType: GrantType, client: &http.Client{}}
}

func (this *HwPush) Auth(clientid, clientSecret string) error {
	req, err := NewRequest()
	if err != nil {
		return err
	}
	resp, err := this.send(req)
	if err != nil {
		return err
	}
	if resp.code != RESP_STATE_OK {
		return errors.New(resp.message)
	}

	return nil
}

func (this *HwPush) send(r *request) (*response, error) {
	r, err := http.NewRequest(http.MethodPost, r.uri, r.body)
	if err != nil {
		return nil, err
	}
	resp, err := this.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

}

func (this *HwPush) Single() {

}

func (this *HwPush) Group() {

}
