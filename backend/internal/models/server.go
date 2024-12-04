package models

import (
	"context"
	"net/http"
	"time"
)

type ServerApi struct {
	server *http.Server
}

func (s *ServerApi) Run(port string, handlers http.Handler) error {
	s.server = &http.Server{
		Addr:         ":" + port,
		Handler:      handlers,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return s.server.ListenAndServe()
}

func (s *ServerApi) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
