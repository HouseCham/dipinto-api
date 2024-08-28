package model

import "time"

type Address struct {
    ID           uint64    `gorm:"primaryKey;autoIncrement"`
    UserID       uint64    `gorm:"not null"`
    User         User      `gorm:"foreignKey:UserID"`
    AddressLine1 string    `gorm:"not null"`
    AddressLine2 string
    City         string    `gorm:"not null"`
    State        string    `gorm:"not null"`
    PostalCode   string    `gorm:"not null"`
    Country      string    `gorm:"not null"`
    CreatedAt    time.Time `gorm:"autoCreateTime"`
    UpdatedAt    time.Time `gorm:"autoUpdateTime"`
    DeletedAt    *time.Time
}