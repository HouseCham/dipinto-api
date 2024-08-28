package model

import "time"

type Product struct {
    ID            uint64    `gorm:"primaryKey;autoIncrement"`
    Name          string    `gorm:"not null"`
    Description   string
    Price         float64   `gorm:"not null;type:numeric(10,2)"`
    StockQuantity int       `gorm:"not null"`
    CategoryID    uint64    `gorm:"not null"`
    Category      Category  `gorm:"foreignKey:CategoryID"`
    CreatedAt     time.Time `gorm:"autoCreateTime"`
    UpdatedAt     time.Time `gorm:"autoUpdateTime"`
    DeletedAt     *time.Time
}

type ProductSize struct {
    ID            uint64    `gorm:"primaryKey;autoIncrement"`
    ProductID     uint64    `gorm:"not null"`
    Product       Product   `gorm:"foreignKey:ProductID"`
    Size          string    `gorm:"not null"`
    Price         float64   `gorm:"not null;type:numeric(10,2)"`
    StockQuantity int       `gorm:"not null"`
    CreatedAt     time.Time `gorm:"autoCreateTime"`
    UpdatedAt     time.Time `gorm:"autoUpdateTime"`
    DeletedAt     *time.Time
}