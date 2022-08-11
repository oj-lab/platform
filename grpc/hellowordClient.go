package main

import (
	"context"
	"flag"
	pb "github.com/OJ-lab/oj-lab-services/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"time"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "192.168.1.153:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("failed to close connection: %v", err)
		}
	}(conn)
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())

	cc := pb.NewJudgerClient(conn)
	stream, err := cc.Judge(context.Background())
	waitcc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitcc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}
			log.Printf(in.Result)
		}
	}()
	for i := 0; i < 20; i++ {
		req := pb.JudgeRequest{
			Language:    "",
			Code:        "",
			Input:       "",
			Output:      "",
			TimeLimit:   0,
			MemoryLimit: 0,
		}
		if err := stream.Send(&req); err != nil {
			log.Fatalf("Failed to send a note: %v", err)
		}
	}
	err = stream.CloseSend()
	if err != nil {
		log.Fatalf("Failed to close stream: %v", err)
	}
	<-waitcc
}
