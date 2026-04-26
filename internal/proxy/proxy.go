package proxy

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type contextKey string

const RetryKey contextKey = "retry_count"

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

func SetRetryInContext(r *http.Request, count int) *http.Request {
	ctx := context.WithValue(r.Context(), RetryKey, count)
	return r.WithContext(ctx)
}

func GetRetryFromContext(r *http.Request) int {
	if count, ok := r.Context().Value(RetryKey).(int); ok {
		return count
	}
	return 0
}
