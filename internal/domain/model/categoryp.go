package model

import "time"

type Category struct {
    ID          uint64    `gorm:"primaryKey;autoIncrement"`
    Name        string    `gorm:"not null"`
    Description string
    CreatedAt   time.Time `gorm:"autoCreateTime"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}