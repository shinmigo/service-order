package rpc

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"goshop/service-order/model/order"
	"goshop/service-order/model/order_address"
	"goshop/service-order/model/order_item"
	"goshop/service-order/model/order_payment"
	"goshop/service-order/pkg/db"
	"goshop/service-order/pkg/grpc/gclient"
	"goshop/service-order/pkg/utils"

	"github.com/shopspring/decimal"

	"github.com/shinmigo/pb/memberpb"

	"github.com/shinmigo/pb/productpb"

	"github.com/shinmigo/pb/basepb"

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
			ConfirmTime: func() string {
				if order.ConfirmTime.IsZero() {
					return ""
				}

				return order.ConfirmTime.Format(utils.TIME_STD_FORMART)
			}(),
			OrderStatus:   order.OrderStatus,
			RefundStatus:  order.RefundStatus,
			ReturnStatus:  order.ReturnStatus,
			UserNote:      order.UserNote,
			OrderItems:    orderItems,
			OrderAddress:  orderAddress,
			OrderPayment:  orderPayment,
			OrderShipment: orderShipment,
			CreatedAt:     order.CreatedAt.Format(utils.TIME_STD_FORMART),
		})
	}

	return &orderpb.ListOrderRes{
		Total:  total,
		Orders: orderDetails,
	}, nil
}

func (o *Order) GetOrderStatus(ctx context.Context, req *orderpb.GetOrderStatusReq) (*orderpb.ListOrderStatusRes, error) {
	var (
		rows   *sql.Rows
		result = make([]*orderpb.ListOrderStatusRes_OrderStatistics, 0, 8)
		err    error
	)
	if rows, err = db.Conn.Model(&order.Order{}).Where("store_id = ?", req.StoreId).
		Select("order_status, count(*) as count").
		Group("order_status").
		Rows(); err != nil {
		return nil, err
	}

	for rows.Next() {
		var row orderpb.ListOrderStatusRes_OrderStatistics
		db.Conn.ScanRows(rows, &row)
		result = append(result, &row)
	}

	return &orderpb.ListOrderStatusRes{
		OrderStatistics: result,
	}, nil
}

func (o *Order) AddOrder(ctx context.Context, req *orderpb.Order) (*basepb.AnyRes, error) {
	var (
		err              error
		productIds       []uint64
		orderData        *order.Order
		orderAddressData *order_address.OrderAddress
		orderPaymentData *order_payment.OrderPayment
		addressDetail    *memberpb.Address
		paymentCode      string
		ok               bool
		listProductResp  *productpb.ListProductRes
	)
	productSpecs := make(map[uint64]map[uint64]*productpb.ProductSpec) //map[ProductId]map[ProductSpecId]Spec
	productsDetail := make(map[uint64]*productpb.ProductDetail)        //map[ProductId]ProductDetail
	specValues := make(map[uint64]map[uint64]map[string]string)        //map[ProductId]map[SpecValueId][string]string
	orderId := utils.GetUniqueId()

	type specDescriptionChildren struct {
		Content     string `json:"content"`
		SpecId      uint64 `json:"spec_id"`
		SpecValueId uint64 `json:"spec_value_id"`
	}

	type specDescription struct {
		Name     string                             `json:"name"`
		SpecId   uint64                             `json:"spec_id"`
		Children map[string]specDescriptionChildren `json:"children"`
	}

	if req.PaymentCode == 0 {
		return nil, fmt.Errorf("未选择支付类型")
	}

	if paymentCode, ok = orderpb.PaymentCode_name[int32(req.PaymentCode)]; !ok {
		return nil, fmt.Errorf("支付类型不存在")
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

	for _, product := range req.Products {
		productIds = append(productIds, product.ProductId)
	}
	if listProductResp, err = gclient.ProductClient.GetProductList(ctx, &productpb.ListProductReq{
		Page:      1,
		PageSize:  100,
		StoreId:   req.StoreId,
		ProductId: productIds,
	}); err != nil {
		return nil, err
	}

	if listProductResp.Total == 0 {
		return nil, fmt.Errorf("商品不存在")
	}

	for _, product := range listProductResp.Products {
		var specDescriptions []*specDescription
		for _, productSpec := range product.Spec {
			if _, ok := productSpecs[product.ProductId]; ok {
				productSpecs[product.ProductId][productSpec.ProductSpecId] = productSpec
			} else {
				buf := make(map[uint64]*productpb.ProductSpec)
				buf[productSpec.ProductSpecId] = productSpec
				productSpecs[product.ProductId] = buf
			}
		}
		productsDetail[product.ProductId] = product
		json.Unmarshal([]byte(product.SpecDescription), &specDescriptions)
		for _, specDescription := range specDescriptions {
			for _, spec := range specDescription.Children {
				content := make(map[string]string)
				content["name"] = specDescription.Name
				content["value"] = spec.Content
				if _, ok := specValues[product.ProductId]; ok {
					specValues[product.ProductId][spec.SpecValueId] = content
				} else {
					buf := make(map[uint64]map[string]string)
					buf[spec.SpecValueId] = content
					specValues[product.ProductId] = buf
				}
			}
		}
	}

	var (
		subtotalDecimal       decimal.Decimal
		grandTotalDecimal     decimal.Decimal
		totalPaidDecimal      decimal.Decimal
		totalQtyOrdered       uint64
		shippingAmountDecimal decimal.Decimal
		discountAmountDecimal decimal.Decimal
	)
	for _, product := range req.Products {
		productDetail := productsDetail[product.ProductId]
		productSpec := productSpecs[product.ProductId][product.ProductSpecId]
		subtotalDecimal = decimal.NewFromFloat(productSpec.Price).Add(subtotalDecimal)
		totalQtyOrdered += product.Qty
		totalPayable, _ := decimal.NewFromFloat(productSpec.Price).Mul(decimal.NewFromFloat(float64(product.Qty))).Float64()
		orderItem := &order_item.OrderItem{
			OrderItemId: utils.GetUniqueId(),
			StoreId:     req.StoreId,
			ParentId:    0,
			OrderId:     orderId,
			Name:        productDetail.Name,
			Sku:         productSpec.Sku, //填规格SKU
			Image: func() string {
				if productSpec.Image != "" {
					return productSpec.Image
				} else {
					return productDetail.Images[0]
				}
			}(), //填规格图片
			ProductId:           product.ProductId,
			Price:               productSpec.Price,    //选择的规格价格
			OldPrice:            productSpec.OldPrice, //选择的规格价格
			CostPrice:           productSpec.CostPrice,
			TotalPayable:        totalPayable, //可能存在精度丢失问题，该库可以解决 github.com/shopspring/decimal
			TotalDiscountAmount: 0,
			QtyOrdered:          product.Qty,
			QtyShipped:          0,
			Weight:              productSpec.Weight,
			Volume:              productSpec.Volume,
			Spec: func() string {
				var specs []*order_item.Spec
				for _, specValueId := range productSpec.SpecValueId {
					spec := &order_item.Spec{
						Name:        specValues[product.ProductId][specValueId]["name"],
						SpecValueId: specValueId,
						Value:       specValues[product.ProductId][specValueId]["value"],
					}
					specs = append(specs, spec)
				}
				specJson, _ := json.Marshal(&specs)

				return string(specJson)
			}(),
		}

		if err = tx.Create(orderItem).Error; err != nil {
			return nil, err
		}
	}
	grandTotalDecimal = subtotalDecimal.Add(shippingAmountDecimal).Sub(discountAmountDecimal)

	grandTotal, _ := grandTotalDecimal.Float64()
	subtotal, _ := subtotalDecimal.Float64()
	shippingAmount, _ := shippingAmountDecimal.Float64()
	totalPaid, _ := totalPaidDecimal.Float64()
	discountAmount, _ := discountAmountDecimal.Float64()
	//订单表
	orderData = &order.Order{
		OrderId:         orderId,
		StoreId:         req.StoreId,
		MemberId:        req.MemberId,
		Subtotal:        subtotal,
		GrandTotal:      grandTotal,
		TotalPaid:       totalPaid,
		TotalQtyOrdered: totalQtyOrdered,
		ShippingAmount:  shippingAmount,
		DiscountAmount:  discountAmount,
		PaymentType:     orderpb.OrderPaymentType_Online,
		PaymentStatus:   orderpb.OrderPaymentStatus_Unpaid,
		PaymentTime:     utils.JSONTime{},
		ShippingStatus:  orderpb.OrderShippingStatus_NotShipped,
		ShippingTime:    utils.JSONTime{},
		Confirm:         orderpb.OrderConfirm_ConfirmNo,
		ConfirmTime:     utils.JSONTime{},
		OrderStatus:     orderpb.OrderStatus_PendingPayment,
		OrderType:       orderpb.OrderType_Normal,
		RefundStatus:    orderpb.OrderRefundStatus_NotRefund,
		ReturnStatus:    orderpb.OrderReturnStatus_NotReturn,
		UserNote:        req.UserNode,
	}
	if err = tx.Create(orderData).Error; err != nil {
		return nil, err
	}

	if addressDetail, err = gclient.Address.GetAddressDetail(ctx, &basepb.GetOneReq{
		Id: req.AddressId,
	}); err != nil {
		return nil, err
	}
	if addressDetail.MemberId != req.MemberId {
		return nil, fmt.Errorf("选择的地址有误")
	}

	orderAddressData = &order_address.OrderAddress{
		OrderAddressId: utils.GetUniqueId(),
		OrderId:        orderId,
		Receiver:       addressDetail.Name,
		Telephone:      addressDetail.Mobile,
		Province:       addressDetail.ProvName,
		City:           addressDetail.CityName,
		Region:         addressDetail.CounName,
		Street:         addressDetail.Address + addressDetail.RoomNumber,
		CreatedBy:      0,
		UpdatedBy:      0,
	}
	if err = tx.Create(orderAddressData).Error; err != nil {
		return nil, err
	}

	orderPaymentData = &order_payment.OrderPayment{
		OrderPaymentId: utils.GetUniqueId(),
		OrderId:        orderId,
		ShippingAmount: shippingAmount,
		AmountPaid:     0,
		AmountOrdered:  grandTotal,
		PaymentName:    order_payment.GetPaymentName(paymentCode),
		PaymentCode:    paymentCode,
	}

	if err = tx.Create(orderPaymentData).Error; err != nil {
		return nil, err
	}

	if err = tx.Commit().Error; err != nil {
		return nil, err
	}

	return &basepb.AnyRes{
		Id:    orderId,
		State: 1,
	}, nil
}
