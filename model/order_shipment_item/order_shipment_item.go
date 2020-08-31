package order_shipment_item

type OrderShipmentItem struct {
	OrderShipmentItemId uint64 `gorm:"PRIMARY_KEY"`
	OrderItemId         uint64
	OrderShipmentId     uint64
	ProductId           uint64
	Name                string
	Sku                 string
	Price               float64
	Qty                 uint64
	Weight              float64
	Volume              float64
	Spec                string
}

func GetTableName() string {
	return "order_shipment_item"
}
