package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/righ/grpc-go-example/goodnightworld/goodnightworld"
	"google.golang.org/grpc"
)

const (
	address = "grpc-server:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	cli, err := c.SayGoodnight(ctx)
	cli.Send(&pb.GoodnightRequest{Name: "righ"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
		return
	}
	for {
		name, err := cli.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Good night: %s", name)
	}
}
