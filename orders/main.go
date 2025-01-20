package main

import (
	"context"
	"log"
	"net"

	common "github.com/Euclid0192/commons"
	"google.golang.org/grpc"
)

var (
	grpcAddr = common.EnvString("GRPC_ADDR", "localhost:3000")
)

func main() {

	grpcServer := grpc.NewServer()
	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	defer l.Close()

	store := NewStore()
	service := NewService(store)
	NewGRPCHandler(grpcServer, service)

	service.CreateOrder(context.Background()) /// empty Context

	log.Println("GRPC Server Started at", grpcAddr)

	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
