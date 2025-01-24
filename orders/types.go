package main

import (
	"context"

	pb "github.com/Euclid0192/commons/api"
)

type OrdersService interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
	ValidateOrders(context.Context, *pb.CreateOrderRequest) ([]*pb.Item, error)
}

type OrdersStore interface {
	Create(context.Context) error
}
