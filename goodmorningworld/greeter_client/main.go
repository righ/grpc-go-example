package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/righ/grpc-go-example/goodmorningworld/goodmorningworld"
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	stream, err := c.SayGoodmorning(ctx)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
		return
	}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		name := scanner.Text()
		if name == "" {
			break
		}
		if err := stream.Send(&pb.GoodmorningRequest{Name: name}); err != nil {
			log.Fatalf("Send failed: %v", err)
		}
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("could not greet: %v", err)
		return
	}
	fmt.Println(reply.GetMessage())
}
