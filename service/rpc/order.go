package rpc

import (
	"context"
	"goshop/service-order/model/order"
	"goshop/service-order/pkg/utils"

	"github.com/shinmigo/pb/orderpb"
)

type Order struct {
}

func NewOrder() *Order {
	return &Order{}
}

func (o *Order) GetOrderList(ctx context.Context, req *orderpb.ListOrderReq) (*orderpb.ListOrderRes, error) {
	var (
		orderDetails = make([]*orderpb.OrderDetail, 0, req.PageSize)
	)
	orders, total, err := order.GetOrders(req)

	if err != nil {
		return nil, err
	}

	for _, order := range orders {
		var (
			orderItems    = make([]*orderpb.OrderDetail_OrderItems, 0, 8)
			orderAddress  *orderpb.OrderDetail_OrderAddress
			orderPayment  *orderpb.OrderDetail_OrderPayment
			orderShipment *orderpb.OrderDetail_OrderShipment
		)

		for _, item := range order.OrderItem {
			orderItems = append(orderItems, &orderpb.OrderDetail_OrderItems{
				OrderItemId:         item.OrderItemId,
				ProductId:           item.ProductId,
				Name:                item.Name,
				Sku:                 item.Sku,
				Image:               item.Image,
				Price:               item.Price,
				OldPrice:            item.OldPrice,
				TotalPayable:        item.TotalPayable,
				TotalDiscountAmount: item.TotalDiscountAmount,
				QtyOrdered:          item.QtyOrdered,
				Weight:              item.Weight,
				Volume:              item.Volume,
				Spec:                item.Spec,
				QtyShipped:          item.QtyShipped,
			})
		}

		if order.OrderPayment != nil {
			orderPayment = &orderpb.OrderDetail_OrderPayment{
				OrderPaymentId: order.OrderPayment.OrderPaymentId,
				PaymentCode:    order.OrderPayment.PaymentCode,
				PaymentName:    order.OrderPayment.PaymentName,
			}
		}

		if order.OrderAddress != nil {
			orderAddress = &orderpb.OrderDetail_OrderAddress{
				OrderAddressId: order.OrderAddress.OrderAddressId,
				Receiver:       order.OrderAddress.Receiver,
				Telephone:      order.OrderAddress.Telephone,
				Province:       order.OrderAddress.Province,
				City:           order.OrderAddress.City,
				Region:         order.OrderAddress.Region,
				Street:         order.OrderAddress.Street,
			}
		}

		if order.OrderShipment != nil {
			orderShipment = &orderpb.OrderDetail_OrderShipment{
				CarrierName:    order.OrderShipment.CarrierName,
				TrackingNumber: order.OrderShipment.TrackingNumber,
			}
		}

		orderDetails = append(orderDetails, &orderpb.OrderDetail{
			OrderId:        order.OrderId,
			StoreId:        order.StoreId,
			MemberId:       order.MemberId,
			OrderType:      order.OrderType,
			Subtotal:       order.Subtotal,
			GrandTotal:     order.GrandTotal,
			TotalPaid:      order.TotalPaid,
			ShippingAmount: order.ShippingAmount,
			DiscountAmount: order.DiscountAmount,
			PaymentType:    order.PaymentType,
			PaymentStatus:  order.PaymentStatus,
			PaymentTime: func() string {
				if order.PaymentTime.IsZero() {
					return ""
				}

				return order.PaymentTime.Format(utils.TIME_STD_FORMART)
			}(),
			ShippingStatus: order.ShippingStatus,
			ShippingTime: func() string {
				if order.ShippingTime.IsZero() {
					return ""
				}

				return order.ShippingTime.Format(utils.TIME_STD_FORMART)
			}(),
			Confirm: order.Confirm,
			ConfigTime: func() string {
				if order.ConfigTime.IsZero() {
					return ""
				}

				return order.ConfigTime.Format(utils.TIME_STD_FORMART)
			}(),
			OrderStatus:   order.OrderStatus,
			RefundStatus:  order.RefundStatus,
			ReturnStatus:  order.ReturnStatus,
			UserNote:      order.UserNote,
			OrderItems:    orderItems,
			OrderAddress:  orderAddress,
			OrderPayment:  orderPayment,
			OrderShipment: orderShipment,
		})
	}

	return &orderpb.ListOrderRes{
		Total:  total,
		Orders: orderDetails,
	}, nil
}
