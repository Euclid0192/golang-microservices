package main

import (
	"errors"
	"net/http"

	common "github.com/Euclid0192/commons"
	pb "github.com/Euclid0192/commons/api"
	"github.com/Euclid0192/order-management-system-gateway/gateway"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	// gateway
	// client pb.OrderServiceClient
	gateway gateway.OrdersGateway
}

func NewHandler(gateway gateway.OrdersGateway) *handler {
	return &handler{gateway: gateway}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/customers/{customerID}/orders", h.HandleCreateOrder)
}

func (h *handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")

	var items []*pb.ItemsWithQuantity

	if err := common.ReadJSON(r, &items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateItems(items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	o, err := h.gateway.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CusomterID: customerID,
		Items:      items,
	})

	/// need to extract error message form gRPC error code
	rStatus := status.Convert(err)

	if rStatus != nil {
		if rStatus.Code() != codes.InvalidArgument {
			common.WriteError(w, http.StatusBadRequest, rStatus.Message())
			return
		}
		common.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	/// Write newly created order back to client
	common.WriteJSON(w, http.StatusOK, o)
}

func validateItems(items []*pb.ItemsWithQuantity) error {
	if len(items) == 0 {
		return common.ErrorNoItem
	}

	for _, i := range items {
		if i.ID == "" {
			return errors.New("Item should have an ID")
		}

		if i.Quantity <= 0 {
			return errors.New("Item must have valid quantity")
		}
	}

	return nil
}
