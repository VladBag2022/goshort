package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/VladBag2022/goshort/internal/misc"
	pb "github.com/VladBag2022/goshort/internal/proto"
)

type StatsResponse struct {
	Urls  int64 `json:"urls"`
	Users int64 `json:"users"`
}

type Entry struct {
	Result string `json:"short_url"`
	Origin string `json:"original_url"`
}

type BatchShortenRequestEntry struct {
	ID     string `json:"correlation_id"`
	Origin string `json:"original_url"`
}

type BatchShortenResponseEntry struct {
	ID     string `json:"correlation_id"`
	Result string `json:"short_url"`
}

func authCookieHelper(s Server, w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie(s.abstractServer.Config.AuthCookieName)

	var token string
	if err == nil {
		token = cookie.Value
	} else if !errors.Is(err, http.ErrNoCookie) {
		return "", err
	}

	userID, err := s.abstractServer.ValidateOrRegister(r.Context(), token)
	if err != nil {
		return "", err
	}

	cookie = &http.Cookie{
		Name:  s.abstractServer.Config.AuthCookieName,
		Value: misc.Sign(s.abstractServer.Config.AuthCookieKey, userID),
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

		response, err := s.abstractServer.Shorten(r.Context(), userID, &pb.ShortenRequest{Origin: content})
		if err != nil {
			pbStatus, ok := status.FromError(err)
			if ok && pbStatus.Code() == codes.InvalidArgument {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		if response.Existed {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusCreated)
		}

		_, err = w.Write([]byte(response.Result))
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

		var request pb.ShortenRequest
		if err = protojson.Unmarshal(body, &request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := s.abstractServer.Shorten(r.Context(), userID, &request)
		if err != nil {
			pbStatus, ok := status.FromError(err)
			if ok && pbStatus.Code() == codes.InvalidArgument {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		responseBytes, err := protojson.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if response.Existed {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusCreated)
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

		s.abstractServer.Delete(userID, &pb.DeleteRequest{UrlIDs: request})
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

		response, err := s.abstractServer.List(r.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(response.GetEntries()) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		entries := make([]Entry, len(response.GetEntries()))
		for i, pbEntry := range response.GetEntries() {
			entries[i].Result = pbEntry.GetResult()
			entries[i].Origin = pbEntry.GetOrigin()
		}

		responseBytes, marshalErr := json.Marshal(entries)
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

func restoreHandler(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, "The id parameter is missing", http.StatusBadRequest)
			return
		}

		response, err := s.abstractServer.Restore(r.Context(), &pb.RestoreRequest{Id: id})
		if err != nil {
			pbStatus, ok := status.FromError(err)
			if ok && pbStatus.Code() == codes.NotFound {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if response.GetDeleted() {
			w.WriteHeader(http.StatusGone)
			return
		}

		w.Header().Set("Location", response.GetOrigin())
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func pingHandler(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := s.abstractServer.Ping(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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

		var requestList []BatchShortenRequestEntry
		if err = json.Unmarshal(body, &requestList); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var pbRequest pb.BatchShortenRequest
		pbRequest.Entries = make([]*pb.BatchShortenRequestEntry, len(requestList))
		for i, e := range requestList {
			pbRequest.GetEntries()[i] = &pb.BatchShortenRequestEntry{
				Id:     e.ID,
				Origin: e.Origin,
			}
		}

		pbResponse, err := s.abstractServer.ShortenBatch(r.Context(), userID, &pbRequest)
		if err != nil {
			pbStatus, ok := status.FromError(err)
			if ok && pbStatus.Code() == codes.InvalidArgument {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		responseList := make([]BatchShortenResponseEntry, len(pbResponse.GetEntries()))
		for i, e := range pbResponse.GetEntries() {
			responseList[i].ID = e.GetId()
			responseList[i].Result = e.GetResult()
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

func statsHandler(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pbStats, err := s.abstractServer.Stats(r.Context(), r.RemoteAddr)
		if err != nil {
			pbStatus, ok := status.FromError(err)
			if ok && pbStatus.Code() == codes.Unauthenticated {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		stats := StatsResponse{
			Urls:  pbStats.GetUrls(),
			Users: pbStats.GetUsers(),
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
