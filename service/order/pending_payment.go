package order

import (
	"context"
	"errors"
	"goshop/service-order/model/order"
	"goshop/service-order/model/order_shipment"

	"github.com/shinmigo/pb/orderpb"
)

//待付款
type PendingPayment struct {
	order *order.Order
}

func (p *PendingPayment) Pay() error {
	//todo: Pay()
	return nil
}

func (p *PendingPayment) Review() error {
	return errors.New("非法操作")
}

func (p *PendingPayment) Ship(ctx context.Context, req *orderpb.Shipment) (*order_shipment.OrderShipment, error) {
	return nil, errors.New("非法操作")
}

func (p *PendingPayment) Receive() error {
	return errors.New("非法操作")
}

func (p *PendingPayment) Complete() error {
	return errors.New("非法操作")
}

func (p *PendingPayment) Comment() error {
	return errors.New("非法操作")
}

func (p *PendingPayment) Cancel() error {
	//todo : cancel()
	return nil
}
