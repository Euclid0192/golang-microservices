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

func (s *service) CreateOrder(ctx context.Context) error {
	return nil
}

func (s *service) ValidateOrders(ctx context.Context, p *pb.CreateOrderRequest) error {
	if len(p.Items) == 0 {
		return common.ErrorNoItem
	}

	mergedItems := mergeItemsWithQuantities(p.Items)

	log.Printf("Merged items: %v", mergedItems)

	return nil
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
