package bouyguessms

import (
	"errors"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type httpClient interface {
	Get(string) (string, error)
	PostForm(string, url.Values) (string, error)
}

type defaultHTTPClient struct {
	client *http.Client
}

func newHTTPClient() (*defaultHTTPClient, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar: jar,
	}

	return &defaultHTTPClient{client}, nil
}

func (c *defaultHTTPClient) Get(url string) (string, error) {
	res, err := c.client.Get(url)
	if err != nil {
		return "", nil
	}

	return handleResponse(res)
}

func (c *defaultHTTPClient) PostForm(url string, data url.Values) (string, error) {
	res, err := c.client.PostForm(url, data)
	if err != nil {
		return "", nil
	}

	return handleResponse(res)
}

func handleResponse(res *http.Response) (body string, err error) {
	if res.StatusCode != http.StatusOK {
		return body, errors.New(res.Status)
	}

	bodyData, err := ioutil.ReadAll(res.Body)
	defer func() {
		if cerr := res.Body.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	if err != nil {
		return body, err
	}

	return string(bodyData), err
}
