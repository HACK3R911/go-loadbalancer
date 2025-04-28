package proxy

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"

	"cloudru/internal/balancer"
)

type contextKey string

const ctxKeyError = contextKey("error")

type Proxy struct {
	balancer *balancer.Balancer
	proxy    *httputil.ReverseProxy
}

func New(b *balancer.Balancer) *Proxy {
	p := &Proxy{balancer: b}
	p.proxy = &httputil.ReverseProxy{
		Director:     p.director,
		ErrorHandler: p.errorHandler,
	}
	return p
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}

func (p *Proxy) director(req *http.Request) {
	backend, err := p.balancer.Next()
	if err != nil {
		ctx := context.WithValue(req.Context(), ctxKeyError, err)
		*req = *req.WithContext(ctx)
		log.Printf("Error selecting backend: %v", err)
		return
	}

	log.Printf("Forwarding to: %s", backend)
	req.URL.Scheme = "http"
	req.URL.Host = backend
}

func (p *Proxy) errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	if ctxErr := r.Context().Value(ctxKeyError); ctxErr != nil {
		log.Printf("Backend error: %v", ctxErr)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	log.Printf("Proxy error: %v", err)
	w.WriteHeader(http.StatusBadGateway)
}
