package httplib

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/vgmdj/utils/logger"
)

const (
	MinRead = 16 * 1024 // 16kb
)

var (
	hc   *Client
	once sync.Once

	DefaultClientConf = &ClientConf{
		Auth: &DefaultAuth{},
	}
)

//ClientConf client config
type ClientConf struct {
	Auth      Authorization
	Timeout   time.Duration
	KeepAlive time.Duration
}

//Client httplib.Client
type Client struct {
	conf      *ClientConf
	httpCli   *http.Client
	dialer    *net.Dialer
	transport http.RoundTripper
	mutex     sync.RWMutex
}

//UniqueClient return the only client
func UniqueClient(conf *ClientConf) *Client {
	once.Do(func() {
		hc = NewClient(conf)
	})

	return hc
}

//NewClient return httplib.Client
func NewClient(conf *ClientConf) *Client {
	if conf == nil {
		conf = DefaultClientConf
	}

	client := &Client{
		conf: conf,
		dialer: &net.Dialer{
			Timeout:   conf.Timeout,
			KeepAlive: conf.KeepAlive,
		},
	}

	client.transport = &http.Transport{
		DialContext:     client.dialer.DialContext,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client.httpCli = &http.Client{
		Transport: client.transport,
	}

	if err := conf.Auth.CheckFormat(); err != nil {
		logger.Error(err.Error())
		panic("must use correct http auth params")
	}

	return client
}

//SetTransport set the transport
func (c *Client) SetTransport(transport *http.Transport) {
	c.transport = transport
	c.httpCli.Transport = transport
}

//SetConfig set config
func (c *Client) SetConfig(conf *ClientConf) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if conf.Auth.CheckFormat() != nil {
		c.conf.Auth = conf.Auth
	}

	if conf.Timeout > 0 {
		c.dialer.Timeout = conf.Timeout
		c.conf.Timeout = conf.Timeout
	}

	if conf.KeepAlive > 0 {
		c.dialer.KeepAlive = conf.KeepAlive
		c.conf.KeepAlive = conf.KeepAlive
	}
}

//Raw sends an HTTP request and use client do
func (c *Client) Raw(method, uri string, body []byte, v interface{}, headers map[string]string) (err error) {
	logger.Info(uri, string(body), headers)

	request, err := c.NewRequest(method, uri, body, headers)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	contentType := headers[ResponseResultContentType]

	data, err := c.Do(request, &contentType)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return respParser(data, contentType, v)

}

//NewRequest return the client new request
func (c *Client) NewRequest(method, uri string, body []byte, headers map[string]string) (request *http.Request, err error) {
	request, err = http.NewRequest(method, uri, bytes.NewReader(body))
	if err != nil {
		logger.Error(err.Error())
		logger.Error("method:", method, "uri:", uri)
		return
	}

	c.conf.Auth.SetAuth(request)

	for k, v := range headers {
		request.Header.Set(k, v)
	}

	return

}

//Do sends an HTTP request and return response data
func (c *Client) Do(request *http.Request, contentType ...*string) (bts []byte, err error) {
	resp, err := c.httpCli.Do(request)
	if err != nil {
		logger.Error(err.Error())
		logger.Error("发送请求错误")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		logger.Error("bad request , status code :", resp.StatusCode)
		bts, _ = ioutil.ReadAll(resp.Body)
		logger.Error("return data :", string(bts))
		return nil, fmt.Errorf("bad request %d", resp.StatusCode)
	}

	bts, err = readAll(resp.Body, MinRead)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	if len(contentType) != 0 {
		*contentType[0] = resp.Header.Get("Content-type")
	}

	return
}
