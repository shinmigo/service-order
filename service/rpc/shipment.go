package rpc

import (
	"context"
	"goshop/service-order/model/order"
	"goshop/service-order/model/order_shipment"
	order2 "goshop/service-order/service/rpc/order"

	"github.com/shinmigo/pb/basepb"
	"github.com/shinmigo/pb/orderpb"
)

type Shipment struct {
}

func NewShipment() *Shipment {
	return &Shipment{}
}

func (s *Shipment) AddShipment(ctx context.Context, req *orderpb.Shipment) (*basepb.AnyRes, error) {
	var (
		ord      *order.Order
		shipment *order_shipment.OrderShipment
		err      error
	)
	if ord, err = order.GetOneByOrderId(req.OrderId); err != nil {
		return nil, err
	}

	orderOperate := order2.NewOperate(ord)
	if shipment, err = orderOperate.Ship(ctx, req); err != nil {
		return nil, err
	}

	return &basepb.AnyRes{
		Id:    shipment.OrderShipmentId,
		State: 1,
	}, nil
}
