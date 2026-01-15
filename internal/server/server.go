package server

import (
	"cloudru/internal/logger"
	"cloudru/internal/proxy"
	"context"
	"golang.org/x/time/rate"
	"net/http"
)

type Server struct {
	httpServer *http.Server
	logger     *logger.Logger
	proxy      *proxy.Proxy
	limiter    *rate.Limiter
}

func New(logger *logger.Logger, p *proxy.Proxy, rl rate.Limit, burst int) *Server {
	s := &Server{logger: logger, proxy: p}
	mux := http.NewServeMux()
	mux.Handle("/", rateLimitMiddleware(p, rl, burst))
	mux.Handle("/health", http.HandlerFunc(s.health))
	return s
}

func (s *Server) Run(port string) error {
	s.httpServer = &http.Server{
		Addr: ":" + port,
	}
	s.logger.Infof("Запуск сервера на порту %s", port)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Завершение работы сервера...")
	return s.httpServer.Shutdown(ctx)
}

func rateLimitMiddleware(next http.Handler, r rate.Limit, burst int) http.HandlerFunc {
	limiter := rate.NewLimiter(r, burst)
	return func(w http.ResponseWriter, req *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, req)
	}
}
func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
