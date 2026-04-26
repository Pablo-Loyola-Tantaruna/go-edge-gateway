package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type GopherProxy struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

func NewGopherProxy(target *url.URL) *GopherProxy {
	return &GopherProxy{
		target: target,
		proxy:  httputil.NewSingleHostReverseProxy(target),
	}
}

func (p *GopherProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}
