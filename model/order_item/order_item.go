package order_item

import (
	"fmt"
	"goshop/service-order/pkg/db"
	"goshop/service-order/pkg/utils"
)

type OrderItem struct {
	OrderItemId         uint64 `gorm:"PRIMARY_KEY"`
	StoreId             uint64
	ParentId            uint64
	OrderId             uint64
	Name                string
	Sku                 string
	Image               string
	ProductId           uint64
	Price               float64
	OldPrice            float64
	CostPrice           float64
	TotalPayable        float64
	TotalDiscountAmount float64
	QtyOrdered          uint64
	QtyShipped          uint64
	Weight              float64
	Volume              float64
	Spec                string
	CreatedAt           utils.JSONTime
	UpdatedAt           utils.JSONTime
}

type Spec struct {
	Name        string `json:"name"`
	SpecValueId uint64 `json:"spec_value_id"`
	Value       string `json:"value"`
}

func GetTableName() string {
	return "order_item"
}

func GetListByOrderId(orderId uint64) ([]*OrderItem, error) {
	var (
		orderItem []*OrderItem
	)

	if err := db.Conn.Where("order_id = ?", orderId).Find(&orderItem).Error; err != nil {
		return nil, fmt.Errorf("err: %v", err)
	}

	return orderItem, nil
}
