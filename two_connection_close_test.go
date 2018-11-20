package main

//go:generate minimock -i net.Conn -o ./ -p main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTwoConnectionCloseHeaders(t *testing.T) {
	req, err := http.NewRequest("GET", "http://google.com", nil)
	if err != nil {
		panic(err)
	}

	req.Close = true

	client := http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				conn := NewConnMock(t).WriteMock.Set(func(p []byte) (int, error) {
					assert.Equal(t, "GET / HTTP/1.1\r\nHost: google.com\r\nUser-Agent: Go-http-client/1.1\r\nConnection: close\r\nAccept-Encoding: gzip\r\n\r\n", string(p))
					return 0, errors.New("unexpected err")
				})

				return conn, nil
			},
		},
	}

	client.Do(req)
}
