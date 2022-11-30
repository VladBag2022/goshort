package grpc

import (
	"context"

	"github.com/VladBag2022/goshort/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (s Server) userIDInterceptor(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	token := valueFromIncomingContext(ctx, server.DefaultAuthCookieName)

	if len(token) == 0 {
		ctx = metadata.NewIncomingContext(ctx, metadata.Pairs(userIDMetadata, ""))
		return handler(ctx, req)
	}

	userID, err := s.abstractServer.Validate(ctx, token)
	if err != nil {
		return nil, err
	}

	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs(userIDMetadata, userID))
	return handler(ctx, req)
}
