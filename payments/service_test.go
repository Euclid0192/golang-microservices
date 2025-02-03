package main

import (
	"context"
	"testing"

	"github.com/Euclid0192/commons/api"
	inmemRegistry "github.com/Euclid0192/commons/discovery/inmem"
	"github.com/Euclid0192/order-management-system-payments/gateway"
	"github.com/Euclid0192/order-management-system-payments/processor/inmem"
)

func TestService(t *testing.T) {
	processor := inmem.NewInmem()
	registry := inmemRegistry.NewRegistry()
	gateway := gateway.NewGRPCGateway(registry)
	svc := NewService(processor, gateway)

	t.Run("should create a payment link", func(t *testing.T) {
		link, err := svc.CreatePayment(context.Background(), &api.Order{})
		/// Error in create payment
		if err != nil {
			t.Errorf("CreatePayment() error = %v, want nil", err)
		}

		/// Check empty payment link
		if link == "" {
			t.Errorf("CreatePayment() link is empty")
		}
	})
}
