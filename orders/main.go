package main

import (
	"context"
	"log"
	"net"
	"time"

	common "github.com/Euclid0192/commons"
	"github.com/Euclid0192/commons/broker"
	"github.com/Euclid0192/commons/discovery"
	"github.com/Euclid0192/commons/discovery/consul"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
)

var (
	serviceName = "orders"
	grpcAddr    = common.EnvString("GRPC_ADDR", "localhost:3000")
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
	amqpUser    = common.EnvString("RABBITMQ_USER", "guest")
	amqpPass    = common.EnvString("RABBITMQ_PASS", "guest")
	amqpHost    = common.EnvString("RABBITMQ_HOST", "localhost")
	amqpPort    = common.EnvString("RABBITMQ_PORT", "5672")
	jaegerAddr  = common.EnvString("JAEGER_ADDR", "localhost:4318")
)

func main() {
	/// Add tracer
	err := common.SetGlobalTracer(context.TODO(), serviceName, jaegerAddr)
	if err != nil {
		log.Fatal("failed to set global tracer")
	}

	// Register new service for every microservice (template for all services)
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		panic(err)
	}

	/// Go routine check health
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

	/// Connect to RabbitMQ
	ch, close := broker.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		close()
		ch.Close()
	}()

	/// End connect to RabbitMQ
	grpcServer := grpc.NewServer()
	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	defer l.Close()

	store := NewStore()
	service := NewService(store)
	serviceWithTelemetry := NewTelemetryMiddleware(service)
	/// later can add any serviceWithSomething -> Decorator pattern
	NewGRPCHandler(grpcServer, serviceWithTelemetry, ch)

	// service.CreateOrder(context.Background()) /// empty Context
	consumer := NewConsumer(serviceWithTelemetry)
	go consumer.Listen(ch)

	log.Println("GRPC Server Started at", grpcAddr)

	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
