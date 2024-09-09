package model

import "time"

type Order struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint64    `gorm:"not null" json:"user_id"`
	AddressID     uint64    `gorm:"not null" json:"address_id"`
	OrderDate     time.Time `gorm:"autoCreateTime" json:"order_date"`
	DeliveryDate  time.Time `json:"delivery_date"`
	Status        string    `gorm:"not null;check:status IN ('pending','shipped','delivered','cancelled')" json:"status"`
	PaymentMethod string    `gorm:"not null;check:payment_method IN ('cash', 'card)" json:"payment_method"`
	TotalAmount   float64   `gorm:"not null;type:numeric(10,2)" json:"total_amount"`
}

type OrderItem struct {
	ID        uint64  `gorm:"primaryKey;autoIncrement"`
	OrderID   uint64  `gorm:"not null"`
	Order     Order   `gorm:"foreignKey:OrderID"`
	ProductID uint64  `gorm:"not null"`
	Product   Product `gorm:"foreignKey:ProductID"`
	Quantity  int     `gorm:"not null"`
	Price     float64 `gorm:"not null;type:numeric(10,2)"`
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