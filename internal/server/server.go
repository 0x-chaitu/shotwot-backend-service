package server

import (
	"context"
	"net/http"
	"shotwot_backend/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + cfg.HTTP.Port,
			Handler:        handler,
			ReadTimeout:    cfg.HTTP.ReadTimeout,
			WriteTimeout:   cfg.HTTP.WriteTimeout,
			MaxHeaderBytes: cfg.HTTP.MaxHeaderMegabytes << 20,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServeTLS("server.crt", "server.key")
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
