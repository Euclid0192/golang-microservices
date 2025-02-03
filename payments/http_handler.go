package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	pb "github.com/Euclid0192/commons/api"
	"github.com/Euclid0192/commons/broker"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/webhook"
)

type PaymentHTTPHandler struct {
	channel *amqp.Channel
}

func NewPaymentHTTPHandler(channel *amqp.Channel) *PaymentHTTPHandler {
	return &PaymentHTTPHandler{channel: channel}
}

func (h *PaymentHTTPHandler) registerRoutes(router *http.ServeMux) {
	router.HandleFunc("/webhook", h.handleCheckoutWebhook)
}

func (h *PaymentHTTPHandler) handleCheckoutWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	fmt.Fprintf(os.Stdout, "Got body: %s\n", body)

	/// get from .env
	event, err := webhook.ConstructEventWithOptions(
		body,
		r.Header.Get("Stripe-Signature"),
		endpointStripeSecret,
		webhook.ConstructEventOptions{IgnoreAPIVersionMismatch: true},
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// or "checkout.session.completed"
	if event.Type == "checkout.session.completed" {
		var cs stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &cs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if cs.PaymentStatus == "paid" {
			log.Printf("Payment for Checkout Session %v successful!", cs.ID)

			/// Publish to OrderPaid queue
			orderID := cs.Metadata["orderID"]
			customerID := cs.Metadata["customerID"]
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			o := &pb.Order{
				ID:          orderID,
				CustomerID:  customerID,
				Status:      "paid",
				PaymentLink: "",
			}
			marshalledOrder, _ := json.Marshal(o)

			h.channel.PublishWithContext(ctx, broker.OrderPaidEvent, "", false, false, amqp.Publishing{
				ContentType:  "application/json",
				Body:         marshalledOrder,
				DeliveryMode: amqp.Persistent,
			})
		}

		log.Println("Message published order.paid")
	}

	w.WriteHeader(http.StatusOK)
}
