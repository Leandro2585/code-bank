package server

import (
	"github.com/leandro2585/code-bank/feature"
	"github.com/leandro2585/code-bank/infra/grpc/service"
	"github.com/leandro2585/code-bank/infra/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type GRPCServer struct {
	ProcessTransactionFeature feature.TransactionFeature
}

func NewGRPCServer GRPCServer {
	return GRPCServer{}
}

func (g GRPCServer) Serve() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("could not listen tcp port")
	}
	transactionService := service.NewTransactionService()
	transactionService.ProcessTransactionFeature = g.ProcessTransactionFeature
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterPaymentServiceServer(grpcServer, transactionService)
	grpcServer.Serve(lis)
}