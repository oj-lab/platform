package impls

import (
	"context"
	"log"

	"github.com/OJ-lab/oj-lab-services/src/service/proto"
)

type GreeterServer struct {
	proto.UnimplementedGreeterServer
}

func (s *GreeterServer) Greeting(ctx context.Context, request *proto.GreetingRequest) (*proto.GreetingResponse, error) {
	log.Printf("Received: %v", request.GetName())
	return &proto.GreetingResponse{Message: "Hello " + request.GetName()}, nil
}
