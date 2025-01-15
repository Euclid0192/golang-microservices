package main

import "context"

type store struct {
	// add MongoDB instance here later

}

func NewStore() *store {
	return &store{}
}

func (s *store) Create(ctx context.Context) error {
	return nil
}
