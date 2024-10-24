package dto

import (
	"encoding/json"

	"github.com/HouseCham/dipinto-api/internal/domain/model"
)

type OrderDetailsDTO struct {
	// Order Information
	ID            uint64  `json:"id"`
	OrderDate     string  `json:"order_date"`
	Status        string  `json:"status" validate:"required,oneof=pending shipped delivered cancelled"`
	TotalAmount   float64 `json:"total_amount" validate:"required"`
	DeliveryDate  string  `json:"delivery_date" validate:"required"`
	PaymentMethod string  `json:"payment_method" validate:"required,oneof=cash card"`
	// Shipping Information
	TrackingID      string  `json:"tracking_id"`
	ShippingCompany string  `json:"shipping_company" validate:"required"`
	DeliveryCost    float64 `json:"delivery_cost" validate:"required"`
	// User Information
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone" validate:"required"`
	// Address Information
	StreetNumber string `json:"street_number" validate:"required"`
	Department   string `json:"department"`
	Neighborhood string `json:"neighborhood" validate:"required"`
	City         string `json:"city" validate:"required"`
	State        string `json:"state" validate:"required"`
	PostalCode   string `json:"postal_code" validate:"required"`
	// Order Items Information
	Items []OrderItemDTO `json:"items" validate:"required"`
}

type OrderItemDTO struct {
	// Order Item Information
	ID       uint64  `json:"id"`
	Quantity int     `json:"quantity" validate:"required"`
	Price    float64 `json:"price" validate:"required"`
	// Product Information
	ImageUrl string          `json:"image_url" validate:"required"`
	Images   json.RawMessage `json:"images" validate:"required"`
	Name     string          `json:"name" validate:"required"`
	// Product Size Information
	Size     string  `json:"size"`
	Discount float64 `json:"discount"`
}

type OrderAddressDTO struct {
	AddressID  uint64        `json:"address_id"`
	NewAddress model.Address `json:"new_address"`
}

type OrderProductsInformationDTO struct {
	Products []ValidOrderProductDTO `json:"products"`
	Coupon   string                 `json:"coupon"`
}

type ValidOrderProductDTO struct {
	ProductID uint64 `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type OrderProductsToApplyDTO struct {
	Products []OrderItemDTO `json:"products"`
	Coupon   *model.Coupon  `json:"coupon"`
	Subtotal float64        `json:"subtotal"`
}

type OrderDbDTO struct {
	AddressID       uint64              `json:"address_id" validate:"required"`
	CouponID        uint64              `json:"coupon_id" validate:"required"`
	TotalAmount     float64          `json:"total_amount" validate:"required"`
	DeliveryCost    float64          `json:"delivery_cost" validate:"required"`
	ShippingCompany string           `json:"shipping_company" validate:"required"`
	Status          string           `json:"status" validate:"required"`
	PaymentMethod   string           `json:"payment_method" validate:"required"`
	Items           []OrderItemDBDTO `json:"items" validate:"required"`
}

type OrderItemDBDTO struct {
	ProductID int     `json:"product_id" validate:"required"`
	Quantity  int     `json:"quantity" validate:"required"`
	Price     float64 `json:"price" validate:"required"`
}
