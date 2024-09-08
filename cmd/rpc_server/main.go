// Deprecated currently
// Keep it for future use
package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/oj-lab/platform/cmd/rpc_server/impls"
	config_module "github.com/oj-lab/platform/modules/config"
	"github.com/oj-lab/platform/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	portProp = "rpc-server.port"
)

var (
	port = config_module.AppConfig().GetInt(portProp)
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
