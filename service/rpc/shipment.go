package rpc

import (
	"context"
	"errors"
	"goshop/service-order/model/order"
	"goshop/service-order/model/order_address"
	"goshop/service-order/model/order_item"
	"goshop/service-order/model/order_shipment"
	"goshop/service-order/model/order_shipment_item"
	"goshop/service-order/pkg/db"

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
		shipment      *order_shipment.OrderShipment
		orderData     *order.Order
		address       *order_address.OrderAddress
		orderItems    []*order_item.OrderItem
		shipmentItems [][]interface{}
		totalWeight   float64
		totalQty      uint64
		err           error
	)

	if order_shipment.ExistShipment(req.OrderId) {
		return nil, errors.New("请勿重复发货")
	}

	if address, err = order_address.GetOneByOrderId(req.OrderId); err != nil {
		return nil, err
	}

	if orderData, err = order.GetOneByOrderId(req.OrderId); err != nil {
		return nil, err
	}

	if orderItems, err = order_item.GetListByOrderId(req.OrderId); err != nil {
		return nil, err
	}

	for _, orderItem := range orderItems {
		totalQty += orderItem.QtyOrdered
		totalWeight += orderItem.Weight
	}

	shipment = &order_shipment.OrderShipment{
		StoreId:        orderData.StoreId,
		OrderId:        orderData.OrderId,
		MemberId:       orderData.MemberId,
		OrderAddressId: address.OrderAddressId,
		TotalWeight:    totalWeight,
		TotalQty:       totalQty,
		UserNote:       orderData.UserNote,
		CarrierId:      req.CarrierId,
		CarrierName:    "xxx",
		TrackingNumber: req.TrackingNumber,
	}

	tx := db.Conn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = tx.Create(shipment).Error; err != nil {
		return nil, err
	}

	for _, orderItem := range orderItems {
		shipmentItems = append(shipmentItems, []interface{}{req.OrderId, orderItem.OrderItemId, shipment.OrderShipmentId,
			orderItem.ProductId, orderItem.Name, orderItem.Sku, orderItem.Price,
			orderItem.QtyOrdered, orderItem.Weight, orderItem.Volume})
	}

	if err = db.BatchInsert(tx, order_shipment_item.GetTableName(),
		[]string{"order_id", "order_item_id", "order_shipment_id", "product_id", "name", "sku", "price", "qty", "weight", "volume"},
		shipmentItems); err != nil {
		return nil, err
	}

	tx.Commit()

	return &basepb.AnyRes{
		Id:    shipment.OrderShipmentId,
		State: 1,
	}, nil
}
