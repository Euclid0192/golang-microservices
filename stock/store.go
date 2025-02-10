package main

import (
	"context"
	"fmt"

	pb "github.com/Euclid0192/commons/api"
)

type Store struct {
	stock map[string]*pb.Item
}

func NewStore() *Store {
	return &Store{
		stock: map[string]*pb.Item{
			"2": {
				ID:       "2",
				Name:     "Potato Chips",
				PriceID:  "price_1QqmZmP3SxpbeGudtiSa4i9P",
				Quantity: 10,
			},
			"1": {
				ID:       "1",
				Name:     "Cheese Burger",
				PriceID:  "price_1QqmZKP3SxpbeGudwP2hHhHl",
				Quantity: 20,
			},
		},
	}
}

func (s *Store) GetItem(ctx context.Context, id string) (*pb.Item, error) {
	for _, item := range s.stock {
		if item.ID == id {
			return item, nil
		}
	}

	return nil, fmt.Errorf("item not found")
}

func (s *Store) GetItems(ctx context.Context, ids []string) ([]*pb.Item, error) {
	var res []*pb.Item
	for _, id := range ids {
		if i, ok := s.stock[id]; ok {
			res = append(res, i)
		}
	}

	return res, nil
}
