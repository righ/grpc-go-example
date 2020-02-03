package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "github.com/righ/grpc-go-example/chat2/chat2"
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

var serial = 1
var channels = make(map[uint64]chan message)

type server struct {
	pb.UnimplementedChatServer
}

func (s *server) CreateChannel(_ context.Context, u *pb.Null) (*pb.User, error) {
	serial++
	channels[uint64(serial)] = make(chan message)
	return &pb.User{Id: uint64(serial)}, nil
}

func (s *server) Talk(srv pb.Chat_TalkServer) error {

	for {
		req, err := srv.Recv()
		if err != nil {
			return err
		}
		name, msg := req.GetName(), req.GetMessage()
		log.Printf("%s>: %s", name, msg)
		now := time.Now()

		for id, c := range channels {
			if id == req.Id {
				continue
			}

			go func(c chan message) { c <- message{name, msg, now} }(c)
		}
	}
}

func (s *server) GetMessages(user *pb.User, srv pb.Chat_GetMessagesServer) error {
	for {
		m := <-channels[user.Id]
		srv.Send(&pb.MessageReply{Name: m.name, Message: m.text})
	}
	return nil
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
