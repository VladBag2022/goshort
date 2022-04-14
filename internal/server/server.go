package server

import (
	"github.com/VladBag2022/goshort/internal/storage"
	"net/http"
)

func New(repository storage.Repository) Server {
	return Server{repository}
}

type Server struct {
	repository storage.Repository
}

func (s *Server) ListenAndServer() {
	http.HandleFunc("/", s.root)

	server := &http.Server{
		Addr: "localhost:8080",
	}
	server.ListenAndServe()
}