package order

import (
	"context"
	"errors"
	"goshop/service-order/model/order"
	"goshop/service-order/model/order_shipment"

	"github.com/shinmigo/pb/orderpb"
)

type PendingComment struct {
	order *order.Order
}

func (c *PendingComment) Pay() error {
	return errors.New("请勿重复付款")
}

func (c *PendingComment) Review() error {
	return errors.New("非法操作")
}

func (c *PendingComment) Ship(ctx context.Context, req *orderpb.Shipment) (*order_shipment.OrderShipment, error) {
	return nil, errors.New("非法操作")
}

func (c *PendingComment) Receive() error {
	return errors.New("非法操作")
}

func (c *PendingComment) Complete() error {
	return errors.New("非法操作")
}

func (c *PendingComment) Comment() error {
	//todo: comment
	return nil
}

func (c *PendingComment) Cancel() error {
	return errors.New("非法操作")
}
