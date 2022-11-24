package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"

	"github.com/VladBag2022/goshort/internal/misc"
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

type ShortenBatchListEntryAPIRequest struct {
	ID     string `json:"correlation_id"`
	Origin string `json:"original_url"`
}

type ShortenBatchListEntryAPIResponse struct {
	ID     string `json:"correlation_id"`
	Result string `json:"short_url"`
}

type StatsResponse struct {
	Urls  int64 `json:"urls"`
	Users int64 `json:"users"`
}

func statsHandler(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(s.config.TrustedSubnet) == 0 {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		_, n, err := net.ParseCIDR(s.config.TrustedSubnet)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !n.Contains(net.ParseIP(r.RemoteAddr)) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		var stats StatsResponse

		stats.Urls, err = s.repository.UrlsCount(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		stats.Users, err = s.repository.UsersCount(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		responseBytes, marshalErr := json.Marshal(&stats)
		if marshalErr != nil {
			http.Error(w, marshalErr.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(responseBytes)
		if err != nil {
			log.Error(err)
		}
	}
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
	} else if !errors.Is(err, http.ErrNoCookie) {
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
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := authCookieHelper(s, w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, err := io.ReadAll(r.Body)
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
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		urlID, inserted, err := s.repository.Shorten(r.Context(), origin)
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

		if inserted {
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusConflict)
		}

		_, err = w.Write([]byte(fmt.Sprintf("%s/%s", s.config.BaseURL, urlID)))
		if err != nil {
			log.Error(err)
		}
	}
}

func shortenAPIHandler(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := authCookieHelper(s, w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, err := io.ReadAll(r.Body)
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
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		urlID, inserted, err := s.repository.Shorten(r.Context(), origin)
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

		if inserted {
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusConflict)
		}

		_, err = w.Write(responseBytes)
		if err != nil {
			log.Error(err)
		}
	}
}

func deleteAPIHandler(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := authCookieHelper(s, w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var request []string
		if err = json.Unmarshal(body, &request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		go func() {
			err = s.repository.Delete(context.Background(), userID, request)
			if err != nil {
				log.Error(err)
			}
		}()
		w.WriteHeader(http.StatusAccepted)
	}
}

func shortenedListAPIHandler(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
				origin, _, restoreErr := s.repository.Restore(r.Context(), urlID)
				if restoreErr != nil {
					http.Error(w, restoreErr.Error(), http.StatusInternalServerError)
					return
				}

				responseList = append(responseList, ShortenedListEntryAPIResponse{
					Result: fmt.Sprintf("%s/%s", s.config.BaseURL, urlID),
					Origin: origin.String(),
				})
			}

			responseBytes, marshalErr := json.Marshal(&responseList)
			if marshalErr != nil {
				http.Error(w, marshalErr.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			_, err = w.Write(responseBytes)
			if err != nil {
				log.Error(err)
			}
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

		origin, deleted, err := s.repository.Restore(r.Context(), id)
		if err != nil {
			var unknownIDErr *storage.UnknownIDError
			if errors.As(err, &unknownIDErr) {
				http.Error(w, fmt.Sprintf("Unknown id: %s", unknownIDErr.ID), http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if deleted {
			w.WriteHeader(http.StatusGone)
			return
		}

		w.Header().Set("Location", origin.String())
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func pingHandler(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.postgres == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := s.postgres.Ping(r.Context()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func badRequestHandler(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Bad request", http.StatusBadRequest)
}

func shortenBatchAPIHandler(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := authCookieHelper(s, w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var requestList []ShortenBatchListEntryAPIRequest
		if err = json.Unmarshal(body, &requestList); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var origins []*url.URL
		for _, request := range requestList {
			if len(request.Origin) == 0 {
				http.Error(w, "URL was not provided", http.StatusBadRequest)
				return
			}

			origin, parseErr := url.Parse(request.Origin)
			if parseErr != nil {
				http.Error(w, parseErr.Error(), http.StatusBadRequest)
				return
			}

			origins = append(origins, origin)
		}

		ids, err := s.repository.ShortenBatch(r.Context(), origins, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var responseList []ShortenBatchListEntryAPIResponse
		for i := 0; i < len(requestList); i++ {
			response := ShortenBatchListEntryAPIResponse{
				ID:     requestList[i].ID,
				Result: fmt.Sprintf("%s/%s", s.config.BaseURL, ids[i]),
			}
			responseList = append(responseList, response)
		}

		responseBytes, err := json.Marshal(&responseList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(responseBytes)
		if err != nil {
			log.Error(err)
		}
	}
}
