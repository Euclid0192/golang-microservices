package main

import (
	"context"
	"testing"

	"github.com/Euclid0192/commons/api"
	"github.com/Euclid0192/order-management-system-payments/processor/inmem"
)

func TestService(t *testing.T) {
	processor := inmem.NewInmem()
	svc := NewService(processor)

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
