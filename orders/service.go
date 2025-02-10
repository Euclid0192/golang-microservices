package main

import (
	"context"

	common "github.com/Euclid0192/commons"
	pb "github.com/Euclid0192/commons/api"
	"github.com/Euclid0192/order-management-system-orders/gateway"
)

type service struct {
	store   OrdersStore
	gateway gateway.StockGateway
}

func NewService(store OrdersStore, gateway gateway.StockGateway) *service {
	return &service{store, gateway}
}

func (s *service) UpdateOrder(ctx context.Context, p *pb.Order) (*pb.Order, error) {
	if err := s.store.Update(ctx, p.ID, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *service) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.Order, error) {
	o, err := s.store.Get(ctx, p.OrderID, p.CustomerID)
	if err != nil {
		return nil, err
	}

	return o.ToProto(), err
}

func (s *service) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest, items []*pb.Item) (*pb.Order, error) {
	// items, err := s.ValidateOrders(ctx, p)
	// if err != nil {
	// 	return nil, err
	// }

	id, err := s.store.Create(ctx, Order{
		CustomerID:  p.CustomerID,
		Status:      "pending",
		Items:       items,
		PaymentLink: "",
	})
	if err != nil {
		return nil, err
	}

	/// Temp
	o := &pb.Order{
		ID:         id.Hex(),
		CustomerID: p.CustomerID,
		Status:     "pending",
		Items:      items,
	}

	return o, err
}

func (s *service) ValidateOrders(ctx context.Context, p *pb.CreateOrderRequest) ([]*pb.Item, error) {
	if len(p.Items) == 0 {
		return nil, common.ErrorNoItem
	}

	mergedItems := mergeItemsWithQuantities(p.Items)

	// log.Printf("Merged items: %v", mergedItems)

	/// validate with stock service
	isInStock, items, err := s.gateway.CheckIfItemIsInStock(ctx, p.CustomerID, mergedItems)
	if err != nil {
		return nil, err
	}

	if !isInStock {
		return items, common.ErrorNoStock
	}

	return items, nil
}

func mergeItemsWithQuantities(items []*pb.ItemsWithQuantity) []*pb.ItemsWithQuantity {
	merged := make([]*pb.ItemsWithQuantity, 0)
	for _, item := range items {
		found := false
		for _, afterMerged := range merged {
			if afterMerged.ID == item.ID {
				afterMerged.Quantity += item.Quantity
				found = true
				break
			}
		}

		if !found {
			merged = append(merged, item)
		}
	}

	return merged
}
