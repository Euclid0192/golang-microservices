// / Kind of Mock object when need to test for 3rd services (in this case Stripe)
package inmem

import pb "github.com/Euclid0192/commons/api"

type Inmem struct{}

func NewInmem() *Inmem {
	return &Inmem{}
}

func (i *Inmem) CreatePaymentLink(*pb.Order) (string, error) {
	return "dummy-link", nil
}
