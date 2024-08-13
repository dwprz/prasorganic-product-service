package server

import (
	"fmt"
	"net"

	"github.com/dwprz/prasorganic-product-service/src/common/log"
	"github.com/dwprz/prasorganic-product-service/src/core/grpc/interceptor"
	pb "github.com/dwprz/prasorganic-proto/protogen/product"
	"google.golang.org/grpc"
)

type Grpc struct {
	port                     string
	server                   *grpc.Server
	productGrpcHandler       pb.ProductServiceServer
	unaryResponseInterceptor *interceptor.UnaryResponse
}

// this main grpc server
func NewGrpc(port string, pgh pb.ProductServiceServer, uri *interceptor.UnaryResponse) *Grpc {
	return &Grpc{
		port:                     port,
		productGrpcHandler:       pgh,
		unaryResponseInterceptor: uri,
	}
}

func (g *Grpc) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", g.port))
	if err != nil {
		log.Logger.Errorf("failed to listen on port %s : %v", g.port, err)
	}

	log.Logger.Infof("grpc run in port: %s", g.port)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			g.unaryResponseInterceptor.Recovery,
			g.unaryResponseInterceptor.Error,
		))

	g.server = grpcServer

	pb.RegisterProductServiceServer(grpcServer, g.productGrpcHandler)

	if err := grpcServer.Serve(listener); err != nil {
		log.Logger.Errorf("failed to serve grpc on port %s : %v", g.port, err)
	}
}

func (g *Grpc) Stop() {
	g.server.Stop()
}
