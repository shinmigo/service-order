package order_address

import (
	"fmt"
	"goshop/service-order/pkg/db"
	"goshop/service-order/pkg/utils"
)

type OrderAddress struct {
	OrderAddressId uint64 `gorm:"PRIMARY_KEY"`
	OrderId        uint64
	AddressId      uint64
	Receiver       string
	Telephone      string
	Province       string
	City           string
	Region         string
	Street         string
	CreatedBy      uint64
	UpdatedBy      uint64
	CreatedAt      utils.JSONTime
	UpdatedAt      utils.JSONTime
	DeletedAt      *utils.JSONTime
}

func GetOneByOrderId(orderId uint64) (*OrderAddress, error) {
	var (
		orderAddress OrderAddress
		err          error
	)

	if orderId == 0 {
		return nil, fmt.Errorf("order id is null")
	}

	if err = db.Conn.Model(OrderAddress{}).Where("order_id = ?", orderId).Find(&orderAddress).Error; err != nil {
		return nil, fmt.Errorf("err: %v", err)
	}

	return &orderAddress, nil
}
