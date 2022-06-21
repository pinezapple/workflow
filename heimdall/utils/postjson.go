package utils

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// BUG: Concurrency post ????
var (
	httpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 50,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: time.Second * 5,
	}
)

// PostJSON post json to url
func PostJSON(ctx context.Context, timeout int64, url string, data []byte) (respData []byte, statusCode int, err error) {
	// set timeout
	tempTimeOut := httpClient.Timeout
	defer func() {
		httpClient.Timeout = tempTimeOut // reset timeout
	}()
	httpClient.Timeout = time.Second * time.Duration(timeout)

	// create request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, -1, fmt.Errorf("Can not create post request, err: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	// send
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, -1, fmt.Errorf("Send post request error: %v", err)
	}
	defer resp.Body.Close()
	// read response body and set status code
	respData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, -1, fmt.Errorf("Can not read post response body, error: %v", err)
	}
	statusCode = resp.StatusCode

	return respData, statusCode, nil
}
