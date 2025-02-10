package main

import (
	"context"
	"log"

	pb "github.com/Euclid0192/commons/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// / Could be get from env
const (
	DbName          = "orders"
	CollectionsName = "orders"
)

/// Temp inmem db
// var orders = make([]*pb.Order, 0)

type store struct {
	// DB layer: MongoDB
	db *mongo.Client
}

func NewStore(db *mongo.Client) *store {
	return &store{db}
}

func (s *store) Create(ctx context.Context, o Order) (primitive.ObjectID, error) {
	// /// Temp
	// id := "42"

	// orders = append(orders, &pb.Order{
	// 	ID:          id,
	// 	CustomerID:  p.CustomerID,
	// 	Status:      "pending",
	// 	Items:       items,
	// 	PaymentLink: "",
	// })

	collection := s.db.Database(DbName).Collection(CollectionsName)
	/// Can add a model/schema layer
	newOrder, err := collection.InsertOne(ctx, o)
	if err != nil {
		log.Fatal("failed to create an order")
	}

	id := newOrder.InsertedID.(primitive.ObjectID)

	return id, nil
}

func (s *store) Get(ctx context.Context, id, customerID string) (*Order, error) {
	collection := s.db.Database(DbName).Collection(CollectionsName)

	var o Order
	err := collection.FindOne(ctx, bson.M{
		"_id":        id,
		"customerID": customerID,
	}).Decode(&o)

	return &o, err
}

func (s *store) Update(ctx context.Context, id string, newOrder *pb.Order) error {
	collection := s.db.Database(DbName).Collection(CollectionsName)

	oId, _ := primitive.ObjectIDFromHex(id)

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": oId},
		bson.M{"$set": bson.M{
			"paymentLink": newOrder.PaymentLink,
			"status":      newOrder.Status,
		}},
	)

	return err
}
