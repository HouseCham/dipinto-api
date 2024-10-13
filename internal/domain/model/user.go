package model

import "time"

// User represents the user model
type User struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"not null" json:"name" validate:"required,lettersAccentsBlank"`
	Email     string    `gorm:"not null;unique" json:"email" validate:"required,email"`
	Password  string    `gorm:"not null" json:"password" validate:"required"`
	Phone     string    `gorm:"not null" json:"phone" validate:"required,phone"`
	Role      string    `gorm:"not null;check:role IN ('customer','admin')" json:"role" validate:"required,role"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time
}

// LoginUser represents the user model for login
type LoginUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Remember bool   `json:"remember"`
}
