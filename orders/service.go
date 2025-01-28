package main

import (
	"context"
	"log"

	common "github.com/Euclid0192/commons"
	pb "github.com/Euclid0192/commons/api"
)

type service struct {
	store OrdersStore
}

func NewService(store OrdersStore) *service {
	return &service{store: store}
}

func (s *service) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.Order, error) {
	return s.store.Get(ctx, p.OrderID, p.CustomerID)
}

func (s *service) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest, items []*pb.Item) (*pb.Order, error) {
	// items, err := s.ValidateOrders(ctx, p)
	// if err != nil {
	// 	return nil, err
	// }

	id, err := s.store.Create(ctx, p, items)
	if err != nil {
		return nil, err
	}

	/// Temp
	o := &pb.Order{
		ID:         id,
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

	log.Printf("Merged items: %v", mergedItems)

	/// validate with stock service later

	///	Temp hard-coded data
	var itemsWithPrice []*pb.Item
	for _, i := range mergedItems {
		itemsWithPrice = append(itemsWithPrice, &pb.Item{
			PriceID:  "price_1Qkc7xP3SxpbeGudlu23Tljm",
			ID:       i.ID,
			Quantity: i.Quantity,
		})
	}

	return itemsWithPrice, nil
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
