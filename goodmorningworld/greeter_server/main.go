package main

import (
	"io"
	"log"
	"net"
	"strings"

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
func (s *server) SayGoodmorning(srv pb.Greeter_SayGoodmorningServer) error {
	names := []string{}
	for {
		req, err := srv.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("Received: %s", req)
		names = append(names, req.GetName())
	}
	message := strings.Join(names[:], ",")
	srv.SendAndClose(&pb.GoodmorningReply{Message: "Good morning " + message + "!"})
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
