package model

import "time"

type Payment struct {
    ID            uint64    `gorm:"primaryKey;autoIncrement"`
    OrderID       uint64    `gorm:"not null"`
    Order         Order     `gorm:"foreignKey:OrderID"`
    PaymentMethod string    `gorm:"not null"`
    PaymentDate   time.Time `gorm:"autoCreateTime"`
    Amount        float64   `gorm:"not null;type:numeric(10,2)"`
}