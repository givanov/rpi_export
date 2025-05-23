package server

import (
	"context"
	"fmt"
	"github.com/givanov/rpi_export/pkg/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
)

type Server interface {
	Start()
	Stop() error
}

type server struct {
	srv *http.Server
	cfg *config.Config
}

func New(cfg *config.Config) Server {
	return &server{
		cfg: cfg,
	}
}

func (s *server) Start() {
	zap.L().Info(fmt.Sprintf("%s starting...", config.AppName))

	serveMux := http.NewServeMux()
	serveMux.Handle("/", &defaultHandler{})

	serveMux.Handle("/metrics", promhttp.Handler())

	srv := http.Server{
		Addr:    s.cfg.Address,
		Handler: serveMux,
	}

	s.srv = &srv

	zap.L().Error("server exited", zap.Error(srv.ListenAndServe()))
}

func (s *server) Stop() error {
	return s.srv.Shutdown(context.TODO())
}

type defaultHandler struct{}

func (d *defaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Here be dragons"))
}
