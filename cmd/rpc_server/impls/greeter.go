package impls

import (
	"context"
	"log"

	"github.com/oj-lab/oj-lab-platform/proto"
)

type GreeterServer struct {
	proto.UnimplementedGreeterServer
}

func (s *GreeterServer) Greeting(ctx context.Context, request *proto.GreetingRequest) (*proto.GreetingResponse, error) {
	log.Printf("Received: %v", request.GetName())
	return &proto.GreetingResponse{Message: "Hello " + request.GetName()}, nil
}
