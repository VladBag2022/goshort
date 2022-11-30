// Package server contains server business logic abstracted from actual transport.
package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/VladBag2022/goshort/internal/proto"
	"github.com/VladBag2022/goshort/internal/storage"
)

type Server struct {
	Repository storage.Repository
	postgres   *storage.PostgresRepository
	Config     *Config
}

func NewServer(repository storage.Repository, postgres *storage.PostgresRepository, config *Config) Server {
	return Server{
		Repository: repository,
		postgres:   postgres,
		Config:     config,
	}
}

func (s Server) Shorten(
	ctx context.Context,
	userID string,
	request *pb.ShortenRequest,
) (response *pb.ShortenResponse, err error) {
	response = &pb.ShortenResponse{}

	origin, err := url.Parse(request.GetOrigin())
	if err != nil {
		return response, status.Errorf(codes.InvalidArgument, "%s", err)
	}

	urlID, inserted, err := s.Repository.Shorten(ctx, origin)
	if err != nil {
		return response, status.Errorf(codes.Internal, "%s", err)
	}
	response.Result = fmt.Sprintf("%s/%s", s.Config.BaseURL, urlID)
	response.Existed = !inserted

	err = s.Repository.Bind(ctx, urlID, userID)
	if err != nil {
		err = status.Errorf(codes.Internal, "%s", err)
	}
	return response, err
}

func (s Server) Delete(userID string, request *pb.DeleteRequest) {
	go func() {
		err := s.Repository.Delete(context.Background(), userID, request.GetUrlIDs())
		if err != nil {
			log.Error(err)
		}
	}()
}

func (s Server) List(ctx context.Context, userID string) (response *pb.Entries, err error) {
	response = &pb.Entries{}

	urlIDs, err := s.Repository.ShortenedList(ctx, userID)
	if err != nil {
		return response, status.Errorf(codes.Internal, "%s", err)
	}

	for _, urlID := range urlIDs {
		origin, _, restoreErr := s.Repository.Restore(ctx, urlID)
		if restoreErr != nil {
			return response, status.Errorf(codes.Internal, "%s", err)
		}

		response.Entries = append(response.Entries, &pb.Entry{
			Result: fmt.Sprintf("%s/%s", s.Config.BaseURL, urlID),
			Origin: origin.String(),
		})
	}

	return response, nil
}

func (s Server) Restore(ctx context.Context, request *pb.RestoreRequest) (response *pb.RestoreResponse, err error) {
	response = &pb.RestoreResponse{}

	origin, deleted, err := s.Repository.Restore(ctx, request.GetId())
	if err != nil {
		var unknownIDErr *storage.UnknownIDError
		if errors.As(err, &unknownIDErr) {
			return response, status.Errorf(codes.NotFound, "%s", err)
		}
		return response, status.Errorf(codes.Internal, "%s", err)
	}
	response.Origin = origin.String()
	response.Deleted = deleted
	return response, err
}

func (s Server) Ping(ctx context.Context) error {
	if s.postgres == nil {
		return status.Error(codes.Internal, "no database")
	}

	if err := s.postgres.Ping(ctx); err != nil {
		return status.Errorf(codes.Internal, "%s", err)
	}
	return nil
}

func (s Server) Stats(ctx context.Context, remoteAddr string) (stats *pb.Stats, err error) {
	stats = &pb.Stats{}

	if len(s.Config.TrustedSubnet) == 0 {
		return stats, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	_, n, err := net.ParseCIDR(s.Config.TrustedSubnet)
	if err != nil {
		return stats, status.Errorf(codes.Internal, "%s", err)
	}

	if !n.Contains(net.ParseIP(remoteAddr)) {
		return stats, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	stats.Urls, err = s.Repository.UrlsCount(ctx)
	if err != nil {
		return stats, status.Errorf(codes.Internal, "%s", err)
	}

	stats.Users, err = s.Repository.UsersCount(ctx)
	if err != nil {
		err = status.Errorf(codes.Internal, "%s", err)
	}
	return stats, err
}
