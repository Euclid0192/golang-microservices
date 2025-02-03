package main

import pb "github.com/Euclid0192/commons/api"

type CreateOrderRequest struct {
	Order         *pb.Order `json:"order"`
	RedirectToURL string    `json:"redirectToURL"`
}
