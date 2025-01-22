package gateway

import (
	"context"

	pb "github.com/Euclid0192/commons/api"
)

type OrdersGateway interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
}
