package services

import (
	"context"

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
func (s *PaymentService) CreatePreference() (*preference.Response, error) {
	request := preference.Request{
		Items: []preference.ItemRequest{
			{
				ID: 		"1234",
				Title:       "My product",
				Description: "My product description",
				PictureURL: "https://www.mercadopago.com/org-img/MP3/home/logomp3.gif",
				CategoryID: "art",
				CurrencyID: "MXN",
				Quantity:    1,
				UnitPrice:   75.76,
			},
		},
		Payer: &preference.PayerRequest{
			Name:    "John",
			Surname: "Doe",
			Email:   "jhon_doe@email.com",
		},
		Shipments: &preference.ShipmentsRequest{
			ReceiverAddress: &preference.ReceiverAddressRequest{
				ZipCode:      "45020",
				StreetName:   "A las Praderas",
				StreetNumber: "4854",
				CountryName: "Mexico",
				StateName: "Jalisco",
				CityName: "Guadalajara",
			},
			Cost: 100,
			Dimensions: "30x30x30,500",
			LocalPickup: false,
			FreeShipping: false,
			ExpressShipment: false,

		},
		Marketplace: "www.dipinto.com.mx",
	}
	
	resource, err := s.client.Create(context.Background(), request)
	if err != nil {
		log.Warnf("Error creating preference: %v", err)
		return nil, err
	}
	return resource, nil
}