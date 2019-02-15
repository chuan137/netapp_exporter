package netapp

import (
	"crypto/tls"
	"io"
	"net/http"
)

type Client struct {
	client *http.Client
	config *ClientConfig
}

type ClientConfig struct {
	url      string
	username string
	password string
}

func NewClient(cfg *ClientConfig) (c *Client) {
	c.client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	return c
}

func (c *Client) NewRequest(xmlbody io.Reader) (*http.Request, error) {
	r, e := http.NewRequest("POST", c.config.url, xmlbody)
	if e != nil {
		return nil, e
	}

	r.SetBasicAuth(c.config.username, c.config.password)
	r.Header.Add("Content-Type", "text/xml")

	return r, nil
}

// func (c *Client) Do() {}
