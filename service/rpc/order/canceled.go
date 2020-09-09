package order

import (
	"context"
	"errors"
	"goshop/service-order/model/order"
	"goshop/service-order/model/order_shipment"

	"github.com/shinmigo/pb/orderpb"
)

type Canceled struct {
	order *order.Order
}

func (c *Canceled) Pay() error {
	return errors.New("订单已取消，请勿付款")
}

func (c *Canceled) Review() error {
	return errors.New("订单已取消，无法审核")
}

func (c *Canceled) Ship(ctx context.Context, req *orderpb.Shipment) (*order_shipment.OrderShipment, error) {
	return nil, errors.New("订单已取消，无法发货")
}

func (c *Canceled) Receive() error {
	return errors.New("非法操作")
}

func (c *Canceled) Complete() error {
	return errors.New("非法操作")
}

func (c *Canceled) Comment() error {
	return errors.New("订单已取消，无法评论")
}

func (c *Canceled) Cancel() error {
	return errors.New("请勿重复取消")
}
