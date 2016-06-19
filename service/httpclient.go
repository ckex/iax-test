package service

import (
	"net/http"
	"time"
	"net"
	"io/ioutil"
	"net/url"
	"iax-test/models"
)

type IaxRequest struct {
	Id    int
	SubId int
	Descr string
	Url   string
}

type IaxResponse struct {
	IaxRequest
	UseTime float32 // ms
	StatusCode int
	Status  string // e.g. "200 OK"
	Body    string
	Err     error
	Result models.Result
}

var transport = &http.Transport{
	ResponseHeaderTimeout:time.Millisecond * time.Duration(500),
	MaxIdleConnsPerHost:100,
	DisableKeepAlives:true,
	Dial:(&net.Dialer{
		Timeout:60 * time.Second,
		KeepAlive:2 * time.Minute,
	}).Dial,
}

var client = &http.Client{
	Transport:transport,
}

func Get(iaxRequest IaxRequest) IaxResponse {

	res := new(IaxResponse)
	res.IaxRequest = iaxRequest

	u, err := url.QueryUnescape(iaxRequest.Url)
	if err != nil {
		res.Err = err
		return *res
	}
	resp, err := client.Get(u)
	if err != nil {
		res.Err = err
		return *res
	}
	defer resp.Body.Close()
	res.Status = resp.Status
	res.StatusCode = resp.StatusCode
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		res.Err = err
		return *res
	}
	res.Body = string(body)
	return *res
}