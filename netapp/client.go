package netapp

import (
	"bytes"
	"crypto/tls"
	"net/http"
)

type Client struct {
	client   *http.Client
	url      string
	username string
	password string
}

func NewClient() (c *Client) {
	c.client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	return c
}

func (c *Client) NewRequest(xmlbody *bytes.Buffer) (*http.Request, error) {
	r, e := http.NewRequest("POST", c.url, xmlbody)
	if e != nil {
		return nil, e
	}

	r.SetBasicAuth(c.username, c.password)
	r.Header.Add("Content-Type", "text/xml")

	return r, nil
}

// func (c *Client) Do() {}
