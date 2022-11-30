package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/VladBag2022/goshort/internal/misc"
	pb "github.com/VladBag2022/goshort/internal/proto"
)

func userIDFromIncomingContext(ctx context.Context) (string, error) {
	return notNilValueFromIncomingContext(ctx, userIDMetadata)
}

func notNilValueFromIncomingContext(ctx context.Context, name string) (string, error) {
	value := valueFromIncomingContext(ctx, name)

	if len(value) == 0 {
		return "", status.Errorf(codes.Internal, "missing %s", name)
	}

	return value, nil
}

func valueFromIncomingContext(ctx context.Context, name string) (value string) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get(name)
		if len(values) > 0 {
			value = values[0]
		}
	}
	return
}

func (s *Server) Shorten(ctx context.Context, r *pb.ShortenRequest) (*pb.ShortenResponse, error) {
	userID, err := userIDFromIncomingContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.abstractServer.Shorten(ctx, userID, r)
}

func (s *Server) Delete(ctx context.Context, r *pb.DeleteRequest) (*empty.Empty, error) {
	userID, err := userIDFromIncomingContext(ctx)
	if err != nil {
		return nil, err
	}

	s.abstractServer.Delete(userID, r)
	return nil, nil
}

func (s *Server) List(ctx context.Context, _ *empty.Empty) (*pb.Entries, error) {
	userID, err := userIDFromIncomingContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.abstractServer.List(ctx, userID)
}

func (s *Server) Restore(ctx context.Context, r *pb.RestoreRequest) (*pb.RestoreResponse, error) {
	return s.abstractServer.Restore(ctx, r)
}

func (s *Server) Ping(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	return nil, s.abstractServer.Ping(ctx)
}

func (s *Server) ShortenBatch(ctx context.Context, r *pb.BatchShortenRequest) (*pb.BatchShortenResponse, error) {
	userID, err := userIDFromIncomingContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.abstractServer.ShortenBatch(ctx, userID, r)
}

func (s *Server) Register(ctx context.Context, _ *empty.Empty) (*pb.RegisterResponse, error) {
	userID, err := s.abstractServer.Register(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		Token: misc.Sign(s.abstractServer.Config.AuthCookieKey, userID),
	}, nil
}
