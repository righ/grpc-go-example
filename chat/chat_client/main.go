package main

import (
	"bufio"
	"context"
	"flag"
	"io"
	"log"
	"os"

	pb "github.com/righ/grpc-go-example/chat/chat"
	"google.golang.org/grpc"
)

const (
	address = "grpc-server:50051"
)

func main() {
	var (
		name = flag.String("name", "noname", "user name")
	)
	flag.Parse()

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewChatClient(conn)

	ctx := context.Background()

	stream, err := c.Talk(ctx)
	if err != nil {
		log.Fatalf("could not talk: %v", err)
		return
	}

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			message := scanner.Text()
			if message == "" {
				break
			}
			if err := stream.Send(&pb.MessageRequest{Name: *name, Message: message}); err != nil {
				log.Fatalf("Send failed: %v", err)
			}
		}
	}()
	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("could not chat: %v", err)
		}
		log.Printf("%s> %s", reply.Name, reply.Message)
	}
}
