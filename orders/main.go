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
	"github.com/Euclid0192/order-management-system-orders/gateway"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
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
	/// Global logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	/// Add tracer
	err := common.SetGlobalTracer(context.TODO(), serviceName, jaegerAddr)
	if err != nil {
		// log.Fatal("failed to set global tracer")
		logger.Fatal("could not set global tracer", zap.Error(err))
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
				logger.Error("failed to do health check", zap.Error(err))
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

	/// gateway to call other services
	gateway := gateway.NewGateway(registry)

	store := NewStore()
	service := NewService(store, gateway)
	serviceWithTelemetry := NewTelemetryMiddleware(service)
	/// later can add any serviceWithSomething -> Decorator pattern
	serviceWithLogging := NewLoggingMiddleware(serviceWithTelemetry)
	NewGRPCHandler(grpcServer, serviceWithLogging, ch)

	// service.CreateOrder(context.Background()) /// empty Context
	consumer := NewConsumer(serviceWithLogging)
	go consumer.Listen(ch)

	// log.Println("GRPC Server Started at", grpcAddr)
	logger.Info("Starting GRPC Server at ", zap.String("port", grpcAddr))

	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
