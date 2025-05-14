package conn

import (
	"net"
	"net/http"
	"time"
)

const (
	DefaultMaxIdleConnsPerHost = 5
)

func newHTTPClient(timeout time.Duration, maxConnsPerHost int) *http.Client {
	return &http.Client{
		Timeout:   timeout,
		Transport: toTransport(timeout, maxConnsPerHost),
	}
}

func toTransport(timeout time.Duration, maxConnsPerHost int) http.RoundTripper {
	return &http.Transport{
		DialContext:         (&net.Dialer{Timeout: timeout, KeepAlive: time.Minute}).DialContext,
		TLSHandshakeTimeout: timeout,
		MaxConnsPerHost:     maxConnsPerHost,
	}
}

func toSecond(duration time.Duration) time.Duration {
	return duration * time.Second
}
