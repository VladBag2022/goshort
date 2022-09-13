package server

import (
	"fmt"
	"net/http"

	"github.com/VladBag2022/goshort/internal/storage"
)

type Server struct {
	repository storage.Repository
	postgres   *storage.PostgresRepository
	config     *Config
}

func NewServer(repository storage.Repository, postgres *storage.PostgresRepository, config *Config) Server {
	return Server{
		repository: repository,
		postgres:   postgres,
		config:     config,
	}
}

func (s Server) ListenAndServer() {
	if err := http.ListenAndServe(s.config.Address, router(s)); err != nil {
		fmt.Println(err)
		return
	}
}
