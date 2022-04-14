package server

import (
	"fmt"
	"github.com/VladBag2022/goshort/internal/storage"
	"log"
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
	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf("%s:%d", s.host, s.port),
			router(s)))
}