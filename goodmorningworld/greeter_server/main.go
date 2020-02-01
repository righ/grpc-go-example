package main

import (
	"context"
	"log"
	"net"

	pb "github.com/righ/grpc-go-example/goodmorningworld/goodmorningworld"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement goodmorningworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayGoodmorning implements goodmorningworld.GreeterServer
func (s *server) SayGoodmorning(ctx context.Context, in *pb.GoodmorningRequest) (*pb.GoodmorningReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.GoodmorningReply{Message: "Good morning " + in.GetName()}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
