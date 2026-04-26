package models

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend struct {
	URL          *url.URL
	Alive        bool
	Mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

func (b *Backend) SetAlive(alive bool) {
	b.Mux.Lock()
	b.Alive = alive
	b.Mux.Unlock()
}

func (b *Backend) IsAlive() bool {
	b.Mux.RLock()
	alive := b.Alive
	b.Mux.RUnlock()
	return alive
}
