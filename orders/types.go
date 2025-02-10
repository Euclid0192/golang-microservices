package main

import (
	"context"

	pb "github.com/Euclid0192/commons/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrdersService interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest, []*pb.Item) (*pb.Order, error)
	GetOrder(context.Context, *pb.GetOrderRequest) (*pb.Order, error)
	ValidateOrders(context.Context, *pb.CreateOrderRequest) ([]*pb.Item, error)
	UpdateOrder(context.Context, *pb.Order) (*pb.Order, error)
}

type OrdersStore interface {
	Create(context.Context, Order) (primitive.ObjectID, error)
	Get(ctx context.Context, id, customerID string) (*Order, error)
	Update(ctx context.Context, id string, p *pb.Order) error
}

// / Model layer for MongoDB
type Order struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	CustomerID  string             `bson:"customerID,omitempty"`
	Status      string             `bson:"status,omitempty"`
	PaymentLink string             `bson:"paymentLink,omitempty"`
	Items       []*pb.Item         `bson:"items,omitempty"`
}

func (o *Order) ToProto() *pb.Order {
	return &pb.Order{
		ID:          o.ID.Hex(),
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		PaymentLink: o.PaymentLink,
	}
}
