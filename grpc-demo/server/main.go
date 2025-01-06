package main

import (
	"context"
	"log"
	"net"

	pb "grpc-demo/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Message: "你好, " + req.Name,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("监听失败: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("服务器启动在 :50051")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("服务失败: %v", err)
	}
}
