package gateway

import (
	"context"

	pb "github.com/Euclid0192/commons/api"
)

type KitchenGateway interface {
	UpdateOrder(context.Context, *pb.Order) error
}
