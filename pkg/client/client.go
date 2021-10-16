package client

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const apiVersion = "1"

type HoneyCombClient struct {
	apiUrl    string
	dataSet   string
	writeKey  string
	userAgent string
}

func NewClient(apiUrl string, dataSet string, writeKey string, agent string) HoneyCombClient {
	return HoneyCombClient{
		apiUrl:    apiUrl,
		dataSet:   dataSet,
		writeKey:  writeKey,
		userAgent: agent,
	}
}

func (c *HoneyCombClient) NewRequest(method string, body interface{}) (*http.Request, error) {
	requestUrl, err := url.Parse(c.apiUrl)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to parse URL %s", c.apiUrl)
		return nil, errors.New(errMsg)
	}
	requestUrl.Path = fmt.Sprintf("/%s/markers/%s", apiVersion, c.dataSet)

	req := &http.Request{}
	if method == "GET" {
		req, err = http.NewRequest(method, requestUrl.String(), nil)
		if err != nil {
			return nil, err
		}
	} else {
		req, err = http.NewRequest(method, requestUrl.String(), body.(io.Reader))
		if err != nil {
			return nil, err
		}
	}

	req.Header.Add("X-Honeycomb-Team", c.writeKey)
	req.Header.Add("X-Honeycomb-Dataset", c.dataSet)
	req.Header.Set("User-Agent", c.userAgent)

	return req, nil
}

func MakeRequest(request *http.Request) (*http.Response, error) {
	var httpClient = http.Client{}
	return httpClient.Do(request)
}

func ReadResponse(response *http.Response) ([]byte, error) {
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf(
			"Failed with %d and message: %s",
			response.StatusCode,
			body,
		)
		return nil, errors.New(errMsg)
	}

	return body, nil
}
