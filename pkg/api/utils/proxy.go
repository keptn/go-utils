package api

import (
	"bytes"
	context "context"
	"io/ioutil"
	"net/http"
)

type ProxyV1Interface interface {
	Proxy(http.ResponseWriter, *http.Request)
}

type ProxyHost struct {
	Host   string
	Scheme string
}

type ProxyHandler struct {
	ProxyHost  ProxyHost
	HttpClient *http.Client
}

func createProxyHandler(proxyHost ProxyHost, httpClient *http.Client) *ProxyHandler {
	return &ProxyHandler{
		ProxyHost:  proxyHost,
		HttpClient: httpClient,
	}
}

func (p *ProxyHandler) Proxy(rw http.ResponseWriter, req *http.Request) {

	fwReq := req.Clone(context.TODO())
	fwReq.RequestURI = ""
	var b bytes.Buffer
	b.ReadFrom(req.Body)
	req.Body = ioutil.NopCloser(&b)
	fwReq.Body = ioutil.NopCloser(bytes.NewReader(b.Bytes()))
	fwReq.URL.Host = p.ProxyHost.Host
	fwReq.URL.Scheme = p.ProxyHost.Scheme

	resp, err := p.HttpClient.Do(fwReq)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for name, headers := range resp.Header {
		for _, h := range headers {
			rw.Header().Set(name, h)
		}
	}

	rw.WriteHeader(resp.StatusCode)

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := rw.Write(respBytes); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
