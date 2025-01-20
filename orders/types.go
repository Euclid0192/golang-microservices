package main

import (
	"context"

	pb "github.com/Euclid0192/commons/api"
)

type OrdersService interface {
	CreateOrder(context.Context) error
	ValidateOrders(context.Context, *pb.CreateOrderRequest) error
}

type OrdersStore interface {
	Create(context.Context) error
}
