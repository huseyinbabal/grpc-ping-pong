package main

import (
	"flag"
	"fmt"
	"github.com/huseyinbabal/grpc-ping-pong/pingpong"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"io"
	"log"
	"net"
)

var port = flag.Int("port", 50051, "Server port")

type pingPongServer struct {
	pingpong.UnimplementedPingPongServiceServer
}

func newPingPongServer() *pingPongServer {
	return &pingPongServer{}
}

func (s *pingPongServer) Ping(stream pingpong.PingPongService_PingServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		if err := stream.Send(&pingpong.PongResponse{Message: in.Message}); err != nil {
			return err
		}
	}
}

func main() {
	flag.Parse()

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	healthServer := health.NewServer()
	healthServer.SetServingStatus("grpc.health.v1.pingpong", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	pingpong.RegisterPingPongServiceServer(grpcServer, newPingPongServer())

	listen, listenErr := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if listenErr != nil {
		log.Fatalf("Failed to listen %v", listenErr)
	}

	err := grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
