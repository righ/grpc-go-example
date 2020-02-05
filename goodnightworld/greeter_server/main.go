package main

import (
	"log"
	"net"

	pb "github.com/righ/grpc-go-example/goodnightworld/goodnightworld"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement goodnightworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayGoodnight implements goodnightworld.GreeterServer
func (s *server) SayGoodnight(in *pb.GoodnightRequest, srv pb.Greeter_SayGoodnightServer) error {
	log.Printf("Received: %s", in.GetName())
	srv.Send(&pb.GoodnightReply{Message: "Good night " + in.GetName() + "!"})
	srv.Send(&pb.GoodnightReply{Message: "Good night " + in.GetName() + "!"})
	srv.Send(&pb.GoodnightReply{Message: "Good night " + in.GetName() + "!"})
	return nil
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
