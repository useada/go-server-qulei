package httpclient

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var (
	httpClient *http.Client
)

type Config struct {
	MaxIdleConns        int `toml:"max_idle_conns"`
	MaxIdleConnsPerHost int `toml:"max_idle_conns_per_host"`
	IdleConnTimeout     int `toml:"idle_conn_timeout"`
}

func Init(conf Config) {
	httpClient = &http.Client{
		Transport: &http.Transport{
			IdleConnTimeout: time.Duration(conf.IdleConnTimeout) * time.Second,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        conf.MaxIdleConns,
			MaxIdleConnsPerHost: conf.MaxIdleConnsPerHost,
		},
		Timeout: 20 * time.Second,
	}
}

type Requester struct {
	*http.Request
}

func NewRequester(method, host string, data []byte) (*Requester, error) {
	req, err := http.NewRequest(method, host, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	return &Requester{Request: req}, nil
}

func (r *Requester) SetHeaders(headers map[string]string) {
	for key, val := range headers {
		r.Header.Set(key, val)
	}
}

func (r *Requester) Do(ctx context.Context, timeout time.Duration) ([]byte, error) {
	if httpClient == nil {
		return nil, errors.New("http client need init")
	}

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resp, err := httpClient.Do(r.Request.WithContext(c))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, errors.New("incorrect http status code")
	}

	return ioutil.ReadAll(resp.Body)
}
