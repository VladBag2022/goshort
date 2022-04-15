package server

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"

	"github.com/VladBag2022/goshort/internal/storage"
)

func shortenHandler(s *Server) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		content := string(body)
		if len(content) == 0 {
			http.Error(w, "Post data should be null", http.StatusBadRequest)
			return
		}

		origin, err := url.Parse(content)
		if err != nil{
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := s.repository.Shorten(r.Context(), origin)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("http://%s:%d/%s", s.host, s.port, id)))
	}
}

func restoreHandler(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, "The id parameter is missing", http.StatusBadRequest)
			return
		}

		origin, err := s.repository.Restore(r.Context(), id)
		if err != nil {
			if _, ok := err.(*storage.UnknownIDError); ok {
				http.Error(w, "Unknown id", http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Location", origin.String())
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func badRequestHandler(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Bad request", http.StatusBadRequest)
}