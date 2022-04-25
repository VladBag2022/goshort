package server

import (
	"fmt"
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
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", s.config.Host, s.config.Port), router(s))
	if err != nil {
		return 
	}
}