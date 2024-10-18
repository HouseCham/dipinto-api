package dto

import "github.com/HouseCham/dipinto-api/internal/domain/model"

type OrderDetailsDTO struct {
	// Order Information
	ID            uint64  `json:"id"`
	OrderDate     string  `json:"order_date"`
	Status        string  `json:"status"`
	TotalAmount   float64 `json:"total_amount"`
	DeliveryDate  string  `json:"delivery_date"`
	PaymentMethod string  `json:"payment_method"`
	// Shipping Information
	TrackingID      string  `json:"tracking_id"`
	ShippingCompany string  `json:"shipping_company"`
	DeliveryCost    float64 `json:"delivery_cost"`
	// User Information
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	// Address Information
	StreetNumber string `json:"street_number"`
	Department   string `json:"department"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	// Order Items Information
	Items []OrderItemDTO `gorm:"-" json:"items"`
}

type OrderItemDTO struct {
	// Order Item Information
	ID       uint64  `json:"id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	// Product Information
	ImageUrl string `json:"image_url"`
	Name     string `json:"name"`
	// Product Size Information
	Size     string  `json:"size"`
	Discount float64 `json:"discount"`
}

type OrderAddressDTO struct {
	AddressID  uint64        `json:"address_id"`
	NewAddress model.Address `json:"new_address"`
}
