package bouyguessms

import (
	"net/url"
)

type fakeClient struct {
	calls      int
	bodies     []string
	errors     []error
	postValues []url.Values
}

func newFakeClient(body string, err error) *fakeClient {
	return &fakeClient{
		calls:  0,
		bodies: []string{body},
		errors: []error{err},
	}
}

func (c *fakeClient) Get(url string) (string, error) {
	body := c.bodies[c.calls]
	err := c.errors[c.calls]
	c.calls++
	return body, err
}

func (c *fakeClient) PostForm(url string, data url.Values) (string, error) {
	body := c.bodies[c.calls]
	err := c.errors[c.calls]
	c.calls++
	c.postValues = append(c.postValues, data)
	return body, err
}
