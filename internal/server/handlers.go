package server

import (
	"github.com/VladBag2022/goshort/internal/storage/errors/unknown_id"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

func (s *Server) root(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.shorten(w, r)

	case http.MethodGet:
		s.restore(w, r)

	default:
		http.Error(w, "Unsupported HTTP method", http.StatusBadRequest)
		return
	}
}

func (s *Server) shorten(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content := string(body)
	origin, err := url.Parse(content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL, err := s.repository.Shorten(r.Context(), *origin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL.Result.String()))
}

func (s *Server) restore(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query().Get("id")
	if param == "" {
		http.Error(w, "The id parameter is missing", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(param)
	if err != nil {
		http.Error(w, "The id parameter is not an integer", http.StatusBadRequest)
		return
	}

	shortURL, err := s.repository.Restore(r.Context(), id)
	if err != nil {
		if _, ok := err.(*unknown_id.UnknownIDError); ok {
			http.Error(w, "Unknown id", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", shortURL.Origin.String())
	w.WriteHeader(http.StatusTemporaryRedirect)
}