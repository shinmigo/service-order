package order_payment

import "goshop/service-order/pkg/utils"

type OrderPayment struct {
	OrderPaymentId uint64 `gorm:"PRIMARY_KEY"`
	OrderId        uint64
	ShippingAmount float64
	AmountPaid     float64
	AmountOrdered  float64
	PaymentName    string
	PaymentCode    string
	CreatedAt      utils.JSONTime
	UpdatedAt      utils.JSONTime
}
