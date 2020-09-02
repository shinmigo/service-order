package order

import (
	"context"
	"goshop/service-order/model/order"
	"goshop/service-order/model/order_shipment"

	"github.com/shinmigo/pb/orderpb"
)

type OperateStatus interface {
	Pay() error
	Review() error
	Ship(ctx context.Context, req *orderpb.Shipment) (*order_shipment.OrderShipment, error)
	Receive() error
	Complete() error
	Comment() error
	Cancel() error
}

type Operate struct {
	order  *order.Order
	status OperateStatus
}

func NewOperate(order *order.Order) *Operate {
	var status OperateStatus
	switch order.OrderStatus {
	case orderpb.OrderStatus_PendingPayment:
		status = &PendingPayment{order: order}
	case orderpb.OrderStatus_PendingReview:
		status = &PendingReview{order: order}
	case orderpb.OrderStatus_PendingShipment:
		status = &PendingShipment{order: order}
	case orderpb.OrderStatus_PendingReceiving:
		status = &PendingReceiving{order: order}
	case orderpb.OrderStatus_PendingCompletion:
		status = &Completed{order: order}
	case orderpb.OrderStatus_PendingComment:
		status = &PendingComment{order: order}
	case orderpb.OrderStatus_Canceled:
		status = &Canceled{order: order}
	default:
		status = &PendingPayment{order: order}
	}
	return &Operate{order: order, status: status}
}

func (o *Operate) Pay() error {
	panic("implement me")
}

func (o *Operate) Review() error {
	panic("implement me")
}

func (o *Operate) Ship(ctx context.Context, req *orderpb.Shipment) (*order_shipment.OrderShipment, error) {
	return o.status.Ship(ctx, req)
}

func (o *Operate) Receive() error {
	panic("implement me")
}

func (o *Operate) Complete() error {
	panic("implement me")
}

func (o *Operate) Comment() error {
	panic("implement me")
}

func (o *Operate) Cancel() error {
	panic("implement me")
}
