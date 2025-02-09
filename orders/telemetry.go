package main

import (
	"context"
	"fmt"

	pb "github.com/Euclid0192/commons/api"
	"go.opentelemetry.io/otel/trace"
)

type TelemetryMiddleware struct {
	next OrdersService
}

func NewTelemetryMiddleware(next OrdersService) OrdersService {
	return &TelemetryMiddleware{next}
}

func (s *TelemetryMiddleware) UpdateOrder(ctx context.Context, p *pb.Order) (*pb.Order, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("UpdateOrder: %v", p))
	return s.next.UpdateOrder(ctx, p)
}

func (s *TelemetryMiddleware) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.Order, error) {

	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("GetOrder: %v", p))

	return s.next.GetOrder(ctx, p)
}

func (s *TelemetryMiddleware) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest, items []*pb.Item) (*pb.Order, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("CreateOrder: %v", p))
	return s.next.CreateOrder(ctx, p, items)
}

func (s *TelemetryMiddleware) ValidateOrders(ctx context.Context, p *pb.CreateOrderRequest) ([]*pb.Item, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("ValidateOrder: %v", p))
	return s.next.ValidateOrders(ctx, p)
}
