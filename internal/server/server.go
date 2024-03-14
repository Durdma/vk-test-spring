package server

import (
	"context"
	"net/http"
	"vk-test-spring/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, mux *http.ServeMux) *Server {
	return &Server{httpServer: &http.Server{
		Addr:         ":" + cfg.HTTP.Port,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		Handler:      mux,
	}}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
