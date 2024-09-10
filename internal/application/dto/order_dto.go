package dto

import "encoding/json"

type OrderDetailsDTO struct {
	// Order Information
	ID              uint64  `json:"id"`
	OrderDate       string  `json:"order_date"`
	Status          string  `json:"status"`
	TotalAmount     float64 `json:"total_amount"`
	DeliveryDate    string  `json:"delivery_date"`
	PaymentMethod   string  `json:"payment_method"`
	TrackingID      string  `json:"tracking_id"`
	ShippingCompany string  `json:"shipping_company"`
	// User Information
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	// Address Information
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	// Order Items Information
	Items []OrderItemDTO `gorm:"-",json:"items"`
}

type OrderItemDTO struct {
	// Order Item Information
	ID       uint64  `json:"id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	// Product Information
	Images json.RawMessage `json:"images"`
	Name   string          `json:"name"`
	// Product Size Information
	Size string `json:"size"`
}
