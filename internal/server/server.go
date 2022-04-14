package server

import (
	"fmt"
	"github.com/VladBag2022/goshort/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", router(s))

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", s.host, s.port), r))
}