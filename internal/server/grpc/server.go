package grpc

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	pb "github.com/VladBag2022/goshort/internal/proto"
	"github.com/VladBag2022/goshort/internal/server"
)

const (
	userIDMetadata     = "userID"
	remoteAddrMetadata = "remoteAddr"
)

type Server struct {
	pb.UnimplementedShortenerServer

	abstractServer *server.Server
	g              *grpc.Server
}

func NewServer(abstractServer *server.Server) Server {
	return Server{
		abstractServer: abstractServer,
	}
}

func (s Server) ListenAndServe() {
	listen, err := net.Listen("tcp", s.abstractServer.Config.GRPCAddress)
	if err != nil {
		fmt.Println(err)
		return
	}

	s.g = grpc.NewServer(grpc.UnaryInterceptor(s.userIDInterceptor))
	pb.RegisterShortenerServer(s.g, &s)

	if err = s.g.Serve(listen); err != nil {
		fmt.Println(err)
	}
}

func (s Server) Shutdown() {
	s.g.GracefulStop()
}
