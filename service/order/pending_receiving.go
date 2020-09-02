package order

import (
	"context"
	"errors"
	"goshop/service-order/model/order"
	"goshop/service-order/model/order_shipment"

	"github.com/shinmigo/pb/orderpb"
)

type PendingReceiving struct {
	order *order.Order
}

func (p *PendingReceiving) Pay() error {
	return errors.New("请勿重复付款")
}

func (p *PendingReceiving) Review() error {
	return errors.New("非法操作")
}

func (p *PendingReceiving) Ship(ctx context.Context, req *orderpb.Shipment) (*order_shipment.OrderShipment, error) {
	return nil, errors.New("非法操作")
}

func (p *PendingReceiving) Receive() error {
	//todo: receive
	return nil
}

func (p *PendingReceiving) Complete() error {
	return errors.New("非法操作")
}

func (p *PendingReceiving) Comment() error {
	return errors.New("非法评论，请先确认收货")
}

func (p *PendingReceiving) Cancel() error {
	return errors.New("已发货，无法取消")
}
