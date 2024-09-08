package model

import "time"

type Category struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"not null" validate:"required,lettersAccentsBlank" json:"name"`
	Description string    `validate:"lettersAccentsBlank" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
