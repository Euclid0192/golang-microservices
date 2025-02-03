package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	common "github.com/Euclid0192/commons"
	"github.com/Euclid0192/commons/broker"
	"github.com/Euclid0192/commons/discovery"
	"github.com/Euclid0192/commons/discovery/consul"
	"github.com/Euclid0192/order-management-system-payments/gateway"
	stripeProcessor "github.com/Euclid0192/order-management-system-payments/processor/stripe"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stripe/stripe-go/v78"
	"google.golang.org/grpc"
)

var (
	serviceName          = "payment"
	grpcAddr             = common.EnvString("GRPC_ADDR", "localhost:3001")
	httpAddr             = common.EnvString("HTTP_ADDR", "localhost:8081")
	consulAddr           = common.EnvString("CONSUL_ADDR", "localhost:8500")
	amqpUser             = common.EnvString("RABBITMQ_USER", "guest")
	amqpPass             = common.EnvString("RABBITMQ_PASS", "guest")
	amqpHost             = common.EnvString("RABBITMQ_HOST", "localhost")
	amqpPort             = common.EnvString("RABBITMQ_PORT", "5672")
	stripeKey            = common.EnvString("STRIPE_KEY", "")
	endpointStripeSecret = common.EnvString("STRIPE_ENDPOINT_SECRET", "whsec...")
)

func main() {
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

	/// Stripe setup
	stripe.Key = stripeKey

	/// Connect to RabbitMQ
	ch, close := broker.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		close()
		ch.Close()
	}()

	/// End connect to RabbitMQ

	/// Consumer to consume messages
	stripeProcessor := stripeProcessor.NewProcessor()
	gateway := gateway.NewGRPCGateway(registry)
	service := NewService(stripeProcessor, gateway)
	amqpConsumer := NewConsumer(service)

	go amqpConsumer.Listen(ch)

	/// http server
	mux := http.NewServeMux()
	httpServer := NewPaymentHTTPHandler(ch)
	httpServer.registerRoutes(mux)

	go func() {
		log.Printf("Starting http server at %s", httpAddr)
		if err := http.ListenAndServe(httpAddr, mux); err != nil {
			log.Fatal("failed to start http server")
		}
	}()

	/// Start gRPC server
	grpcServer := grpc.NewServer()
	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	defer l.Close()

	log.Println("GRPC Server Started at", grpcAddr)

	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
