package main

import (
	"context"

	pb "github.com/Euclid0192/commons/api"
	"github.com/Euclid0192/order-management-system-payments/gateway"
	"github.com/Euclid0192/order-management-system-payments/processor"
)

type service struct {
	processor processor.PaymentProcessor
	gateway   gateway.OrdersGateway
}

func NewService(processor processor.PaymentProcessor, gateway gateway.OrdersGateway) *service {
	return &service{processor, gateway}
}

func (s *service) CreatePayment(ctx context.Context, p *pb.Order) (string, error) {
	/// Connect to a payment processor
	paymentLink, err := s.processor.CreatePaymentLink(p)
	if err != nil {
		return "", err
	}

	/// update order with the link later
	err = s.gateway.UpdateOrderAfterPaymentLink(ctx, p.ID, paymentLink)
	if err != nil {
		return "", nil
	}
	return paymentLink, nil
	// return "some link here...", nil
}
