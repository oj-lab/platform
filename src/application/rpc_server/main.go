package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/OJ-lab/oj-lab-services/src/application/rpc_server/impls"
	"github.com/OJ-lab/oj-lab-services/src/core"
	"github.com/OJ-lab/oj-lab-services/src/service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	portProp = "rpc-server.port"
)

var (
	port = core.AppConfig.GetInt(portProp)
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterGreeterServer(s, &impls.GreeterServer{})
	proto.RegisterStreamerServer(s, &impls.StreamerServer{})

	reflection.Register(s)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
