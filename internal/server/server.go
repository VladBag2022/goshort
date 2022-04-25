package server

import (
	"net/http"

	"github.com/VladBag2022/goshort/internal/storage"
)

type Server struct {
	repository storage.Repository
	config     Config
}

func NewServer(repository storage.Repository, config Config) Server {
	return Server{
		repository: repository,
		config:     config,
	}
}

func (s *Server) ListenAndServer() {
	err := http.ListenAndServe(s.config.Address, router(s))
	if err != nil {
		return 
	}
}