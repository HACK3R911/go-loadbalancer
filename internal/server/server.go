package server

import (
	"log"
	"net/http"
)

type Server struct {
	handler http.Handler
	port    string
}

func New(handler http.Handler, port string) *Server {
	return &Server{
		handler: loggingMiddleware(handler),
		port:    port,
	}
}

func (s *Server) Start() error {
	addr := ":" + s.port
	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, s.handler)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
