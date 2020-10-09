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

const OrderPayment_Wechat = "微信"
const OrderPayment_Alipay = "支付宝"

//获取
func GetPaymentName(paymentCode string) (paymentName string) {
	switch paymentCode {
	case "Wechat":
		paymentName = OrderPayment_Wechat
	case "Alipay":
		paymentName = OrderPayment_Alipay
	default:
		paymentName = "未知"
	}

	return
}
