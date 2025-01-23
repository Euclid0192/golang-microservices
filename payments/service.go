package main

import (
	"context"

	pb "github.com/Euclid0192/commons/api"
)

type service struct {
}

func NewService() *service {
	return &service{}
}

func (s *service) CreatePayment(ctx context.Context, p *pb.Order) (string, error) {
	/// Connect to a payment processor

	return "stripe link...", nil
}
