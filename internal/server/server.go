package server

import (
	"cloudru/internal/logger"
	"context"
	"net/http"
)

type Server struct {
	httpServer *http.Server
	logger     *logger.Logger
}

func New(logger *logger.Logger) *Server {
	return &Server{logger: logger}
}

func (s *Server) Run(port string /*, _ http.Handler*/) error {
	s.httpServer = &http.Server{
		Addr: ":" + port,
		//Handler: loggingMiddleware(handler),
	}
	s.logger.Infof("Запуск сервера на порту %s", port)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Завершение работы сервера...")
	return s.httpServer.Shutdown(ctx)
}

//func loggingMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		log.Printf("%s %s", r.Method, r.URL.Path)
//		next.ServeHTTP(w, r)
//	})
//}
