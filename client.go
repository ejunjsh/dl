package main

import (
	"net/http"
)

type client struct {
	_client *http.Client
	url string
}

func newClient(url string) *client{
	return &client{
		&http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		},url,
	}
}
