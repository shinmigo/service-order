package order_shipment

import (
	"goshop/service-order/pkg/db"
	"goshop/service-order/pkg/utils"
)

type OrderShipment struct {
	OrderShipmentId uint64 `gorm:"PRIMARY_KEY"`
	StoreId         uint64
	OrderId         uint64
	MemberId        uint64
	OrderAddressId  uint64
	TotalWeight     float64
	TotalQty        uint64
	UserNote        string
	CarrierId       uint64
	CarrierName     string
	TrackingNumber  string
	CreatedAt       utils.JSONTime
	UpdatedAt       utils.JSONTime
	DeletedAt       *utils.JSONTime
}

func ExistShipment(orderId uint64) bool {
	var (
		shipment OrderShipment
	)
	db.Conn.Select("order_shipment_id").Where("order_id = ?", orderId).First(&shipment)

	return shipment.OrderShipmentId > 0
}
