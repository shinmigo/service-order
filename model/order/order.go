package order

import (
	"fmt"
	"goshop/service-order/model/order_address"
	"goshop/service-order/model/order_item"
	"goshop/service-order/model/order_payment"
	"goshop/service-order/model/order_shipment"
	"goshop/service-order/pkg/db"
	"goshop/service-order/pkg/utils"

	"github.com/jinzhu/gorm"

	"github.com/shinmigo/pb/orderpb"
)

type Order struct {
	OrderId         uint64 `gorm:"PRIMARY_KEY"`
	StoreId         uint64
	MemberId        uint64
	Subtotal        float64
	GrandTotal      float64
	TotalPaid       float64
	TotalQtyOrdered uint64
	ShippingAmount  float64
	DiscountAmount  float64
	PaymentType     orderpb.OrderPaymentType
	PaymentStatus   orderpb.OrderPaymentStatus
	PaymentTime     utils.JSONTime
	ShippingStatus  orderpb.OrderShippingStatus
	ShippingTime    utils.JSONTime
	Confirm         orderpb.OrderConfirm
	ConfigTime      utils.JSONTime
	OrderStatus     orderpb.OrderStatus
	OrderType       orderpb.OrderType
	RefundStatus    orderpb.OrderRefundStatus
	ReturnStatus    orderpb.OrderReturnStatus
	UserNote        string
	CreatedAt       utils.JSONTime
	UpdatedAt       utils.JSONTime
	DeletedAt       *utils.JSONTime
	OrderItem       []*order_item.OrderItem       `gorm:"foreignkey:OrderId"`
	OrderAddress    *order_address.OrderAddress   `gorm:"foreignkey:OrderId"`
	OrderPayment    *order_payment.OrderPayment   `gorm:"foreignkey:OrderId"`
	OrderShipment   *order_shipment.OrderShipment `gorm:"foreignkey:OrderId"`
}

func GetOneByOrderId(orderId uint64) (*Order, error) {
	if orderId == 0 {
		return nil, fmt.Errorf("order id is null")
	}
	row := &Order{}
	err := db.Conn.
		Where("order_id = ?", orderId).
		First(row).Error

	if err != nil {
		return nil, fmt.Errorf("err: %v", err)
	}
	return row, nil
}

func GetOrders(req *orderpb.ListOrderReq) ([]*Order, uint64, error) {
	var (
		orders     = make([]*Order, 0, req.PageSize)
		conditions = make([]func(db *gorm.DB) *gorm.DB, 0, 8)
		err        error
		total      uint64
	)

	query := db.Conn.Model(Order{})

	if req.StoreId > 0 {
		conditions = append(conditions, func(db *gorm.DB) *gorm.DB {
			return db.Where("store_id = ?", req.StoreId)
		})
	}

	if req.MemberId > 0 {
		conditions = append(conditions, func(db *gorm.DB) *gorm.DB {
			return db.Where("member_id = ?", req.MemberId)
		})
	}

	if req.OrderId > 0 {
		conditions = append(conditions, func(db *gorm.DB) *gorm.DB {
			return db.Where("order_id = ?", req.OrderId)
		})
	}

	if req.OrderStatus > 0 {
		conditions = append(conditions, func(db *gorm.DB) *gorm.DB {
			return db.Where("order_status = ?", req.OrderStatus)
		})
	}

	if req.StartCreatedAt != "" {
		conditions = append(conditions, func(db *gorm.DB) *gorm.DB {
			return db.Where("created_at >= ?", req.StartCreatedAt)
		})
	}

	if req.EndCreatedAt != "" {
		conditions = append(conditions, func(db *gorm.DB) *gorm.DB {
			return db.Where("created_at < ?", req.EndCreatedAt)
		})
	}

	if err = query.Preload("OrderItem").
		Preload("OrderAddress").
		Preload("OrderPayment").
		Preload("OrderShipment").
		Scopes(conditions...).
		Order("order_id desc").
		Offset(req.PageSize * (req.Page - 1)).
		Limit(req.PageSize).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	if err = query.Scopes(conditions...).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}
