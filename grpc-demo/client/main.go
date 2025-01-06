package main

import (
	"context"
	"log"
	"time"

	pb "grpc-demo/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "小明"})
	if err != nil {
		log.Fatalf("调用失败: %v", err)
	}
	log.Printf("收到响应: %s", r.GetMessage())
}
