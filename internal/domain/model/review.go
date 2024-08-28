package model

import "time"

type Review struct {
    ID        uint64    `gorm:"primaryKey;autoIncrement"`
    ProductID uint64    `gorm:"not null"`
    Product   Product   `gorm:"foreignKey:ProductID"`
    UserID    uint64    `gorm:"not null"`
    User      User      `gorm:"foreignKey:UserID"`
    Rating    int       `gorm:"not null;check:rating >= 1 AND rating <= 5"`
    Comment   string
    CreatedAt time.Time `gorm:"autoCreateTime"`
}