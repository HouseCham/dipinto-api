package model

import "time"

type Address struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint64    `gorm:"not null" json:"user_id"`
	Alias          string    `json:"alias"`
	Reference      string    `json:"reference"`
	Addressee      string    `json:"addressee"`
	Phone          string    `json:"phone" validate:"phone"`
	AddresseeEmail string    `json:"addressee_email"`
	StreetNumber   string    `gorm:"not null" json:"street_number" validate:"required"`
	Department     string    `json:"department_number"`
	Neighborhood   string    `gorm:"not null" json:"neighborhood" validate:"required"`
	City           string    `gorm:"not null" json:"city" validate:"required"`
	State          string    `gorm:"not null" json:"state" validate:"required"`
	PostalCode     string    `gorm:"not null" json:"postal_code" validate:"required"`
	Country        string    `gorm:"not null" json:"country"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      *time.Time
}
