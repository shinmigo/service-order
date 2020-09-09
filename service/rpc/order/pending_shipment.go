package order

import (
	"context"
	"errors"
	"goshop/service-order/model/order"
	"goshop/service-order/model/order_address"
	"goshop/service-order/model/order_item"
	"goshop/service-order/model/order_shipment"
	"goshop/service-order/model/order_shipment_item"
	"goshop/service-order/pkg/db"
	"goshop/service-order/pkg/grpc/gclient"

	"github.com/shinmigo/pb/shoppb"

	"github.com/shinmigo/pb/orderpb"
)

//待发货
type PendingShipment struct {
	order *order.Order
}

func (p *PendingShipment) Pay() error {
	return errors.New("请勿重复付款")
}

func (p *PendingShipment) Review() error {
	return errors.New("请勿重复审核")
}

func (p *PendingShipment) Ship(ctx context.Context, req *orderpb.Shipment) (*order_shipment.OrderShipment, error) {
	var (
		shipment      *order_shipment.OrderShipment
		address       *order_address.OrderAddress
		orderItems    []*order_item.OrderItem
		shipmentItems [][]interface{}
		totalWeight   float64
		totalQty      uint64
		listCarriers  *shoppb.ListCarrierRes
		carrierName   string
		err           error
	)

	if order_shipment.ExistShipment(req.OrderId) {
		return nil, errors.New("请勿重复发货")
	}

	if address, err = order_address.GetOneByOrderId(req.OrderId); err != nil {
		return nil, err
	}

	if orderItems, err = order_item.GetListByOrderId(req.OrderId); err != nil {
		return nil, err
	}

	for _, orderItem := range orderItems {
		totalQty += orderItem.QtyOrdered
		totalWeight += orderItem.Weight
	}

	if listCarriers, err = gclient.ShopCarrier.GetCarrierList(ctx, &shoppb.ListCarrierReq{
		Id: req.CarrierId,
	}); err != nil {
		return nil, err
	}

	if listCarriers.Total == 0 {
		return nil, errors.New("选择的物流公司不存在")
	}

	for _, carrier := range listCarriers.Carriers {
		carrierName = carrier.Name
	}

	shipment = &order_shipment.OrderShipment{
		StoreId:        p.order.StoreId,
		OrderId:        p.order.OrderId,
		MemberId:       p.order.MemberId,
		OrderAddressId: address.OrderAddressId,
		TotalWeight:    totalWeight,
		TotalQty:       totalQty,
		UserNote:       p.order.UserNote,
		CarrierId:      req.CarrierId,
		CarrierName:    carrierName,
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

	db.Conn.Model(order.Order{}).Where("order_id = ?", req.OrderId).Update(map[string]interface{}{
		"shipping_status": orderpb.OrderShippingStatus_AllShipped,
		"order_status":    orderpb.OrderStatus_PendingReceiving,
	})

	tx.Commit()

	return shipment, nil
}

func (p *PendingShipment) Receive() error {
	return errors.New("非法操作")
}

func (p *PendingShipment) Complete() error {
	return errors.New("非法操作")
}

func (p *PendingShipment) Comment() error {
	return errors.New("非法操作")
}

func (p *PendingShipment) Cancel() error {
	//todo : cancel()
	return nil
}
