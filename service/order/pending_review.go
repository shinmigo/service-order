package order

import (
	"context"
	"errors"
	"goshop/service-order/model/order"
	"goshop/service-order/model/order_shipment"

	"github.com/shinmigo/pb/orderpb"
)

//待审核
type PendingReview struct {
	order *order.Order
}

func (p *PendingReview) Pay() error {
	return errors.New("请勿重复付款")
}

func (p *PendingReview) Review() error {
	//todo: Review()
	return nil
}

func (p *PendingReview) Ship(ctx context.Context, req *orderpb.Shipment) (*order_shipment.OrderShipment, error) {
	return nil, errors.New("请先审核")
}

func (p *PendingReview) Receive() error {
	return errors.New("非法操作")
}

func (p *PendingReview) Complete() error {
	return errors.New("非法操作")
}

func (p *PendingReview) Comment() error {
	return errors.New("非法操作")
}

func (p *PendingReview) Cancel() error {
	//todo : cancel()
	return nil
}
