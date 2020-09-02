package order

import (
	"context"
	"errors"
	"goshop/service-order/model/order"
	"goshop/service-order/model/order_shipment"

	"github.com/shinmigo/pb/orderpb"
)

type Completed struct {
	order *order.Order
}

func (c *Completed) Pay() error {
	return errors.New("订单已完成，请勿重复付款")
}

func (c *Completed) Review() error {
	return errors.New("非法操作")
}

func (c *Completed) Ship(ctx context.Context, req *orderpb.Shipment) (*order_shipment.OrderShipment, error) {
	return nil, errors.New("订单已完成，请勿重复发货")
}

func (c *Completed) Receive() error {
	return errors.New("非法操作")
}

func (c *Completed) Complete() error {
	return errors.New("非法操作")
}

func (c *Completed) Comment() error {
	return errors.New("订单已完成，无法评论")
}

func (c *Completed) Cancel() error {
	return errors.New("非法操作")
}
