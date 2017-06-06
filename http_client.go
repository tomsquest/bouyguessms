package bouyguessms

import (
	"errors"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type HttpClient interface {
	Get(string) (string, error)
	PostForm(string, url.Values) (string, error)
}

type httpClient struct {
	client *http.Client
}

func NewHttpClient() (*httpClient, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar: jar,
	}

	return &httpClient{client}, nil
}

func (c *httpClient) Get(url string) (string, error) {
	res, err := c.client.Get(url)
	if err != nil {
		return "", nil
	}

	return handleResponse(res)
}

func (c *httpClient) PostForm(url string, data url.Values) (string, error) {
	res, err := c.client.PostForm(url, data)
	if err != nil {
		return "", nil
	}

	return handleResponse(res)
}

func handleResponse(res *http.Response) (string, error) {
	if res.StatusCode != http.StatusOK {
		return "", errors.New(res.Status)
	}

	bodyData, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", err
	}

	return string(bodyData), nil
}
