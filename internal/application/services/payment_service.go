package services

import "fmt"

import (
	"context"

	"github.com/HouseCham/dipinto-api/internal/application/dto"
	"github.com/gofiber/fiber/v3/log"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/preference"
)

// PaymentService is a service to handle payment operations
type PaymentService struct {
	client preference.Client
}

// NewPaymentService creates a new PaymentService
func NewPaymentService(cfg *config.Config) *PaymentService {
	client := preference.NewClient(cfg)
	return &PaymentService{
		client: client,
	}
}

// CreatePreference creates a new preference in MercadoPago API
func (s *PaymentService) CreatePreference(order *dto.OrderDetailsDTO) (*preference.Response, error) {
	request := preference.Request{
		Items: generateItemRequest(order.Items),
		Payer: &preference.PayerRequest{
			Name:  order.Name,
			Email: order.Email,
			Phone: &preference.PhoneRequest{
				Number: order.Phone,
			},
		},
		Shipments: &preference.ShipmentsRequest{
			ReceiverAddress: &preference.ReceiverAddressRequest{
				StreetName:  order.StreetNumber,
				Apartment:   order.Department,
				StateName:   order.State,
				ZipCode:     order.PostalCode,
				CityName:    order.City,
				CountryName: "Mexico",
			},
			Cost:            order.DeliveryCost,
			LocalPickup:     false,
			FreeShipping:    false,
			ExpressShipment: false,
			DefaultShippingMethod: "standard",
		},
		Marketplace: "www.dipinto.com",
	}

	resource, err := s.client.Create(context.Background(), request)
	if err != nil {
		log.Warnf("Error creating preference: %v", err)
		return nil, err
	}
	return resource, nil
}
// generateItemRequest generates a list of ItemRequest from a list of OrderItemDTO
func generateItemRequest(orderItems []dto.OrderItemDTO) []preference.ItemRequest {
	items := make([]preference.ItemRequest, 0)
	for _, item := range orderItems {
		items = append(items, preference.ItemRequest{
			ID:          fmt.Sprint(item.ID),
			Title:       item.Name,
			PictureURL:  item.ImageUrl,
			CurrencyID:  "MXN",
			Quantity:    item.Quantity,
			UnitPrice:   item.Price,
		})
	}
	return items
}