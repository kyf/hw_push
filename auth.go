package hw_push

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (this *HwPush) Auth() error {
	if this.token != nil &&
		this.token.ExpireIn.After(time.Now()) {
		return nil
	}

	data := url.Values{
		"grant_type":    []string{GrantType},
		"client_secret": []string{this.ClientSecret},
		"client_id":     []string{this.ClientId},
	}

	req, err := http.NewRequest(http.MethodPost, AuthUri, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", UserAgent)
	resp, err := this.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(bytes.NewReader(responseData))
	var response authResponse
	err = decoder.Decode(&response)
	if err != nil {
		return err
	}

	if response.Error != 0 {
		return errors.New(fmt.Sprintf("[%d]%s", response.Error, response.ErrorDes))
	}

	duration := time.Duration(response.ExpireIn - ExpireBuff)
	this.tokenLocker.Lock()
	this.token = &AccessToken{Value: response.AccessToken, ExpireIn: time.Now().Add(duration)}
	this.tokenLocker.Unlock()
	return nil
}
