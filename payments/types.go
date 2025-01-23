package main

import (
	"context"

	pb "github.com/Euclid0192/commons/api"
)

type PaymentsService interface {
	CreatePayment(context.Context, *pb.Order) (string, error)
}
