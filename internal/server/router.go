package server

import (
	"github.com/NYTimes/gziphandler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func router(s Server) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(DecompressGZIP)
	r.Use(gziphandler.GzipHandler)

	r.Mount("/debug", middleware.Profiler())
	r.Get("/{id}", restoreHandler(s))
	r.Post("/", shortenHandler(s))
	r.Post("/api/shorten", shortenAPIHandler(s))
	r.Post("/api/shorten/batch", shortenBatchAPIHandler(s))
	r.Get("/api/user/urls", shortenedListAPIHandler(s))
	r.Delete("/api/user/urls", deleteAPIHandler(s))
	r.Get("/api/internal/stats", statsHandler(s))
	r.Get("/ping", pingHandler(s))
	r.MethodNotAllowed(badRequestHandler)
	r.NotFound(badRequestHandler)

	return r
}
