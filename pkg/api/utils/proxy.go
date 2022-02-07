package api

import (
	"bytes"
	"context"
	"fmt"
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
	fmt.Printf("Incoming request: host=%s, path=%s, URL=%s", req.URL.Host, req.URL.Path, req.URL.String())
	fwReq := req.Clone(context.TODO())
	fwReq.RequestURI = ""
	var b bytes.Buffer
	b.ReadFrom(req.Body)
	req.Body = ioutil.NopCloser(&b)
	fwReq.Body = ioutil.NopCloser(bytes.NewReader(b.Bytes()))
	fwReq.URL.Host = p.ProxyHost.Host
	fwReq.URL.Scheme = p.ProxyHost.Scheme

	fmt.Printf("Forwarding Request:  host=%s, path=%s, URL=%s", fwReq.URL.Host, fwReq.URL.Path, fwReq.URL.String())
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

//	var path string
//	if req.URL.RawPath != "" {
//		path = req.URL.RawPath
//	} else {
//		path = req.URL.Path
//	}
//	logger.Debugf("Incoming request: host=%s, path=%s, URL=%s", req.URL.Host, path, req.URL.String())
//	proxyScheme, proxyHost, proxyPath := config.Global.ProxyHost(path)
//
//	if proxyScheme == "" || proxyHost == "" {
//		logger.Error("Could not get proxy Host URL - got empty values of 'proxyScheme' or 'proxyHost'")
//		rw.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//
//	forwardReq, err := http.NewRequest(req. Method, req.URL.String(), req.Body)
//	if err != nil {
//		logger.Errorf("Could not create request to be forwarded: %v", err)
//		rw.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//
//	forwardReq.Header = req.Header
//
//	parsedProxyURL, err := url.Parse(proxyScheme + "://" + strings.TrimSuffix(proxyHost, "/") + "/" + strings.TrimPrefix(proxyPath, "/"))
//	if err != nil {
//		logger.Errorf("Could not decode url with scheme: %s, host: %s, path: %s - %v", proxyScheme, proxyHost, proxyPath, err)
//		rw.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//
//	forwardReq.URL = parsedProxyURL
//	forwardReq.URL.RawQuery = req.URL.RawQuery
//	logger.Debugf("Forwarding request to host=%s, path=%s, URL=%s", proxyHost, proxyPath, forwardReq.URL.String())
//
//	if config.Global.KeptnAPIToken != "" {
//		logger.Debug("Adding x-token header to HTTP request")
//		forwardReq.Header.Add("x-token", config.Global.KeptnAPIToken)
//	}
//
//	client := f.httpClient
//	resp, err := client.Do(forwardReq)
//	if err != nil {
//		logger.Errorf("Could not send request to API endpoint: %v", err)
//		rw.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	defer resp.Body.Close()
//
//	for name, headers := range resp.Header {
//		for _, h := range headers {
//			rw.Header().Set(name, h)
//		}
//	}
//
//	rw.WriteHeader(resp.StatusCode)
//
//	respBytes, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		logger.Errorf("Could not read response payload: %v", err)
//		rw.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//
//	logger.Debugf("Received response from API: Status=%d", resp.StatusCode)
//	if _, err := rw.Write(respBytes); err != nil {
//		logger.Errorf("could not send response from API: %v", err)
//	}
//}
