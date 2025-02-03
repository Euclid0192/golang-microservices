package main

import (
	"context"

	pb "github.com/Euclid0192/commons/api"
)

type OrdersService interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest, []*pb.Item) (*pb.Order, error)
	GetOrder(context.Context, *pb.GetOrderRequest) (*pb.Order, error)
	ValidateOrders(context.Context, *pb.CreateOrderRequest) ([]*pb.Item, error)
	UpdateOrder(context.Context, *pb.Order) (*pb.Order, error)
}

type OrdersStore interface {
	Create(context.Context, *pb.CreateOrderRequest, []*pb.Item) (string, error)
	Get(ctx context.Context, id, customerID string) (*pb.Order, error)
	Update(ctx context.Context, id string, p *pb.Order) error
}
