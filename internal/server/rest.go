package server

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"todo/pkg/logger"
)

type RESTServer struct {
	server *http.Server
	gsTime time.Duration
}

type RESTServerConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	GSTime       time.Duration
}

func NewRESTServer(handler http.Handler, cfg RESTServerConfig) *RESTServer {
	return &RESTServer{
		server: &http.Server{
			Addr:         cfg.Addr,
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		gsTime: cfg.GSTime,
	}
}

func (s *RESTServer) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		logger.Info.Printf("сервер слушает на %s", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error.Printf("ошибка HTTP сервера: %v", err)
		}
	}()

	<-ctx.Done()

	logger.Info.Println("завершение HTTP сервера...")

	ctx, cancel = context.WithTimeout(context.Background(), s.gsTime)
	defer cancel()

	return s.Shutdown(ctx)
}

func (s *RESTServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
