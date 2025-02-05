package main

import (
	"context"
	"log"
	"net/http"
	"time"

	common "github.com/Euclid0192/commons"
	"github.com/Euclid0192/commons/discovery"
	"github.com/Euclid0192/commons/discovery/consul"
	"github.com/Euclid0192/order-management-system-gateway/gateway"
	_ "github.com/joho/godotenv/autoload"
)

var (
	httpAddr         = common.EnvString("HTTP_ADDR", ":3000")
	orderServiceAddr = "localhost:3000"
	consulAddr       = common.EnvString("CONSUL_ADDR", "localhost:8500")
	serviceName      = "gateway"
	jaegerAddr       = common.EnvString("JAEGER_ADDR", "localhost:4318")
)

func main() {

	/// Old: connect directly to grpc server
	// conn, err := grpc.NewClient(orderServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// defer conn.Close()

	// if err != nil {
	// 	log.Fatalf("Failed to dial server: %v", err)
	// }
	// log.Printf("Dialing order service at %s", orderServiceAddr)

	/// New: connect through a gateway

	/// Add tracer
	err := common.SetGlobalTracer(context.TODO(), serviceName, jaegerAddr)
	if err != nil {
		log.Fatal("failed to set global tracer")
	}

	// Register new service for every microservice
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, httpAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatalf("failed to health check")
			}
			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)
	/// End register service

	mux := http.NewServeMux()
	ordersGateway := gateway.NewGRPCGateway(registry)

	handler := NewHandler(ordersGateway)
	handler.registerRoutes(mux)

	log.Printf("Starting http server on port: %s", httpAddr)
	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatalf("Failed to start http server: %v", err)
	}

}
