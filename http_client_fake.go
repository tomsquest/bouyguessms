package bouyguessms

import (
	"net/url"
)

type FakeClient struct {
	calls      int
	bodies     []string
	errors     []error
	postValues []url.Values
}

func NewFakeClient(body string, err error) *FakeClient {
	return &FakeClient{
		calls:  0,
		bodies: []string{body},
		errors: []error{err},
	}
}

func (c *FakeClient) Get(url string) (string, error) {
	body := c.bodies[c.calls]
	err := c.errors[c.calls]
	c.calls++
	return body, err
}

func (c *FakeClient) PostForm(url string, data url.Values) (string, error) {
	body := c.bodies[c.calls]
	err := c.errors[c.calls]
	c.calls++
	c.postValues = append(c.postValues, data)
	return body, err
}
