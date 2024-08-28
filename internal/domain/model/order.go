package model

import "time"

type Order struct {
    ID          uint64    `gorm:"primaryKey;autoIncrement"`
    UserID      uint64    `gorm:"not null"`
    User        User      `gorm:"foreignKey:UserID"`
    OrderDate   time.Time `gorm:"autoCreateTime"`
    Status      string    `gorm:"not null;check:status IN ('pending','shipped','delivered','cancelled')"`
    TotalAmount float64   `gorm:"not null;type:numeric(10,2)"`
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