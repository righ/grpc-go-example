package main

import (
	"bufio"
	"context"
	"flag"
	"io"
	"log"
	"os"

	pb "github.com/righ/grpc-go-example/chat2/chat2"
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

	user, err := c.CreateChannel(ctx, &pb.Null{})
	if err != nil {
		log.Fatalf("could not send: %v", err)
		return
	}

	go func() {
		cli, err := c.GetMessages(context.Background(), &pb.User{Id: user.Id})
		if err != nil {
			log.Fatalf("could not get: %v", err)
			return
		}
		for {
			reply, err := cli.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Recv failed: %v", err)
				return
			}
			log.Printf("%s> %s", reply.Name, reply.Message)
		}
	}()

	cli, err := c.SendMessages(ctx)
	if err != nil {
		log.Fatalf("could not send: %v", err)
		return
	}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		if message == "" {
			continue
		}
		if err := cli.Send(&pb.MessageRequest{Id: user.Id, Name: *name, Message: message}); err != nil {
			log.Fatalf("Send failed: %v", err)
		}
	}
	cli.CloseAndRecv()
}
