package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"

	"workflow/workflow-utils/model"
	"workflow/executor/core"

	retry "github.com/avast/retry-go"
)

const (
//TODO: define endpoint
)

var httpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 50,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	},
}

// PostJSON post json data
func PostJSON(url string, withSignature bool, sig string, data []byte, expect interface{}) (err error) {
	rd := bytes.NewReader(data)
	var body []byte

	// create new request
	req, err := http.NewRequest(http.MethodPost, url, rd)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	//TODO: maybe add signature

	er := retry.Do(func() error {
		// make request
		resp, err := httpClient.Do(req)
		if err != nil {
			return err
		}

		// get response body
		body, err = ioutil.ReadAll(resp.Body)
		defer func() {
			_ = resp.Body.Close()
		}()

		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			tmp := template.JSEscapeString(string(body))
			return fmt.Errorf("%d:%s", resp.StatusCode, tmp)
		}
		// try to parse response
		if expect != nil {
			err = json.Unmarshal(body, expect)
		}

		return nil
	}, retry.Attempts(uint(core.GetMainConfig().APIRetryCount)))

	if er != nil {
		return er
	}

	return
}

func ProcessResponse(resp *model.Response) (data interface{}, err error) {
	if resp.Error.Message != "" {
		return nil, fmt.Errorf(resp.Error.Message)
	} else {
		return resp.Data, nil
	}
}
