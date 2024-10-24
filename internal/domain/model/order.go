package model

import "time"

// Order represents an order in the orders table.
type Order struct {
	ID              uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          uint64     `gorm:"not null" json:"user_id"`
	AddressID       uint64     `gorm:"not null" json:"address_id"`
	CouponID        uint64     `gorm:"index" json:"coupon_id,omitempty"`
	OrderDate       time.Time  `gorm:"autoCreateTime" json:"order_date"`
	DeliveryDate    *time.Time `gorm:"index" json:"delivery_date,omitempty"`
	TotalAmount     float64    `gorm:"type:numeric(10,2);not null" json:"total_amount"`
	TrackingID      string     `gorm:"type:text" json:"tracking_id,omitempty"`
	DeliveryCost    float64    `gorm:"type:numeric(10,2)" json:"delivery_cost,omitempty"`
	ShippingCompany string     `gorm:"type:text" json:"shipping_company,omitempty"`
	Status          string     `gorm:"type:text;check:status IN ('pending', 'shipped', 'delivered', 'cancelled');not null" json:"status"`
	PaymentMethod   string     `gorm:"type:text;check:payment_method IN ('cash', 'card');not null" json:"payment_method"`
}

// OrderItem represents an item in an order in the order_items table.
type OrderItem struct {
	ID        uint64   `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID   uint64   `gorm:"not null" json:"order_id"`
	ProductID uint64   `gorm:"not null" json:"product_id"`
	Quantity  int      `gorm:"not null" json:"quantity"`
	Price     float64  `gorm:"type:numeric(10,2);not null" json:"price"`
	Discount  *float64 `gorm:"type:numeric" json:"discount,omitempty"`
}

type OrderListItem struct {
	ID              uint64    `json:"id"`
	OrderDate       time.Time `json:"order_date"`
	Name            string    `json:"name"`
	PaymentMethod   string    `json:"payment_method"`
	TotalAmount     float64   `json:"total_amount"`
	Status          string    `json:"status"`
	DeliveryDate    time.Time `json:"delivery_date"`
	TrackingID      string    `json:"tracking_id"`
	ShippingCompany string    `json:"shipping_company"`
}
