package server

import (
	"fmt"
	"github.com/VladBag2022/goshort/internal/storage"
	"net/http"
)

func New(repository storage.Repository, host string, port int) Server {
	return Server{
		repository: repository,
		host: 		host,
		port: 		port,
	}
}

type Server struct {
	repository 	storage.Repository
	host 		string
	port 		int
}

func (s *Server) ListenAndServer() {
	http.HandleFunc("/", s.root)

	server := &http.Server{
		Addr: fmt.Sprintf("%s:%d", s.host, s.port),
	}
	server.ListenAndServe()
}