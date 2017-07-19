package krakenapi

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type ApiClient struct {
	Key       string
	secret    string
	ApiRoot   string
	UserAgent string
	client    *http.Client
}

func NewApiClient(api_root, key, secret string) *ApiClient {
	client := &http.Client{}

	return &ApiClient{key, secret, api_root, "", client}
}

func (api *ApiClient) Query(url_path string, params url.Values, with_signature bool) ([]byte, error) {
	headers := map[string]string{}
	method := "GET"

	if with_signature {
		secret, _ := base64.StdEncoding.DecodeString(api.secret)
		params.Set("nonce", fmt.Sprintf("%d", time.Now().UnixNano()))

		signature := createKrakenSignature(url_path, params, secret)

		headers["API-Key"] = api.Key
		headers["API-Sign"] = signature

		method = "POST"
	}

	if api.UserAgent != "" {
		headers["User-Agent"] = api.UserAgent
	}

	headers["Content-Type"] = "application/x-www-form-urlencoded"

	return executeHttpQuery(method, api.ApiRoot+url_path, headers, params)
}
