package stripe

import (
	"fmt"
	"log"

	common "github.com/Euclid0192/commons"
	pb "github.com/Euclid0192/commons/api"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
)

var (
	gatewayHTTPAddr = common.EnvString("GATEWAY_HTTP_ADDR", "http://localhost:8080")
)

type Stripe struct {
}

func NewProcessor() *Stripe {
	return &Stripe{}
}

func (s *Stripe) CreatePaymentLink(p *pb.Order) (string, error) {
	log.Printf("Creating payment link for order %v", p)
	/// Gateway succesful URL
	gatewaySuccessURL := fmt.Sprintf("%s/success.html?customerID=%s&orderID=%s", gatewayHTTPAddr, p.CustomerID, p.ID)
	gatewayCancelURL := fmt.Sprintf("%s/cancel.html", gatewayHTTPAddr)

	/// Create params from order
	items := []*stripe.CheckoutSessionLineItemParams{}
	for _, item := range p.Items {
		items = append(items, &stripe.CheckoutSessionLineItemParams{
			Price: stripe.String(item.PriceID),
			// Currently hardcode to test first
			// Price:    stripe.String("price_1Qkc7xP3SxpbeGudlu23Tljm"),
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}

	params := &stripe.CheckoutSessionParams{
		Metadata: map[string]string{
			"orderID":    p.ID,
			"customerID": p.CustomerID,
		},
		LineItems:  items,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(gatewaySuccessURL),
		CancelURL:  stripe.String(gatewayCancelURL),
	}
	result, err := session.New(params)
	if err != nil {
		return "", err
	}
	return result.URL, nil
}
