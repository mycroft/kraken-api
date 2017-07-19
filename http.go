package krakenapi

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func executeHttpQuery(method string, url string, headers map[string]string, values url.Values) ([]byte, error) {
	var bodyReader io.Reader

	client := &http.Client{}

	if method == "GET" {
		bodyReader = nil
		url = url + "?" + values.Encode()
	} else {
		bodyReader = strings.NewReader(values.Encode())
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("Could not execute request! (%s)", err.Error())
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Could not execute request! (%s)", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not execute request! (%s)", err.Error())
	}

	return body, nil
}
