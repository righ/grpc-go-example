package main

import (
	"log"
	"net"
	"time"

	pb "github.com/righ/grpc-go-example/chat/chat"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type message struct {
	name string
	text string
	time time.Time
}

var messages = []message{}

type server struct {
	pb.UnimplementedChatServer
}

func (s *server) Talk(srv pb.Chat_TalkServer) error {
	lastRead := time.Now()

	for {
		req, err := srv.Recv()
		if err != nil {
			return err
		}
		name, msg := req.GetName(), req.GetMessage()
		log.Printf("%s>: %s", name, msg)

		for _, m := range messages {
			if m.time.After(lastRead) {
				srv.Send(&pb.MessageReply{Name: m.name, Message: m.text})
			}
		}
		lastRead = time.Now()
		messages = append(messages, message{name, msg, lastRead})
	}
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterChatServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
