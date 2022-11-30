// Package http contains HTTP API server.
package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VladBag2022/goshort/internal/server"
)

type Server struct {
	abstractServer     *server.Server
	http       *http.Server
}

func NewServer(abstractServer *server.Server) Server {
	s := Server{
		abstractServer: abstractServer,
	}
	s.http = &http.Server{ //nolint:gosec // do not support timeouts for simplicity
		Addr:    abstractServer.Config.Address,
		Handler: router(s),
	}
	return s
}

func (s Server) ListenAndServe() {
	var err error
	if s.abstractServer.Config.EnableHTTPS {
		err = s.http.ListenAndServeTLS(s.abstractServer.Config.CertPEMFile,
			s.abstractServer.Config.KeyPEMFile)
	} else {
		err = s.http.ListenAndServe()
	}
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (s Server) Shutdown() error {
	return s.http.Shutdown(context.Background())
}
