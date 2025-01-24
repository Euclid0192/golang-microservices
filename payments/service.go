package main

import (
	"context"

	pb "github.com/Euclid0192/commons/api"
	"github.com/Euclid0192/order-management-system-payments/processor"
)

type service struct {
	processor processor.PaymentProcessor
}

func NewService(processor processor.PaymentProcessor) *service {
	return &service{processor: processor}
}

func (s *service) CreatePayment(ctx context.Context, p *pb.Order) (string, error) {
	/// Connect to a payment processor
	paymentLink, err := s.processor.CreatePaymentLink(p)
	if err != nil {
		return "", err
	}

	/// update order with the link later

	return paymentLink, nil
	// return "some link here...", nil
}
