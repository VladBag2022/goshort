package server

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/VladBag2022/goshort/internal/misc"
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"

	"github.com/VladBag2022/goshort/internal/storage"
)

type ShortenAPIRequest struct {
	Origin string `json:"url"`
}

type ShortenAPIResponse struct {
	Result string `json:"result"`
}

type ShortenedListEntryAPIResponse struct {
	Result string `json:"short_url"`
	Origin string `json:"original_url"`
}

func authCookieHelper(s Server, w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie(s.config.AuthCookieName)

	if err == nil {
		validCookie, userID, _ := misc.Verify(s.config.AuthCookieKey, cookie.Value)

		if validCookie {
			_, err = s.repository.ShortenedList(r.Context(), userID)
			if err == nil {
				return userID, nil
			}
		}

	} else if err != http.ErrNoCookie {
		return "", err
	}

	userID, err := s.repository.Register(r.Context())
	if err != nil {
		return "", err
	}
	cookie = &http.Cookie{
		Name:  s.config.AuthCookieName,
		Value: misc.Sign(s.config.AuthCookieKey, userID),
	}
	http.SetCookie(w, cookie)

	return userID, nil
}

func shortenHandler(s Server) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		userID, err := authCookieHelper(s, w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var reader io.Reader
		if r.Header.Get(`Content-Encoding`) == `gzip` {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			reader = gz
			defer gz.Close()
		} else {
			reader = r.Body
		}

		body, err := io.ReadAll(reader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		content := string(body)
		if len(content) == 0 {
			http.Error(w, "Post data should not be null", http.StatusBadRequest)
			return
		}

		origin, err := url.Parse(content)
		if err != nil{
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		urlID, err := s.repository.Shorten(r.Context(), origin)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = s.repository.Bind(r.Context(), urlID, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("%s/%s", s.config.BaseURL, urlID)))
	}
}

func shortenAPIHandler(s Server) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		userID, err := authCookieHelper(s, w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var reader io.Reader
		if r.Header.Get(`Content-Encoding`) == `gzip` {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			reader = gz
			defer gz.Close()
		} else {
			reader = r.Body
		}

		body, err := io.ReadAll(reader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var request ShortenAPIRequest
		if err = json.Unmarshal(body, &request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if len(request.Origin) == 0 {
			http.Error(w, "URL was not provided", http.StatusBadRequest)
			return
		}

		origin, err := url.Parse(request.Origin)
		if err != nil{
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		urlID, err := s.repository.Shorten(r.Context(), origin)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = s.repository.Bind(r.Context(), urlID, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := ShortenAPIResponse{
			Result: fmt.Sprintf("%s/%s", s.config.BaseURL, urlID),
		}
		responseBytes, err := json.Marshal(&response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(responseBytes)
	}
}

func shortenedListAPIHandler(s Server) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		userID, err := authCookieHelper(s, w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		urlIDs, err := s.repository.ShortenedList(r.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(urlIDs) == 0 {
			w.WriteHeader(http.StatusNoContent)
		} else {
			var responseList []ShortenedListEntryAPIResponse

			for _, urlID := range urlIDs {
				origin, err := s.repository.Restore(r.Context(), urlID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				responseList = append(responseList, ShortenedListEntryAPIResponse{
					Result: fmt.Sprintf("%s/%s", s.config.BaseURL, urlID),
					Origin: origin.String(),
				})
			}

			responseBytes, err := json.Marshal(&responseList)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(responseBytes)
		}
	}
}

func restoreHandler(s Server) http.HandlerFunc {
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