package model

import (
	"time"
)

type Coupon struct {
	ID            int64     `json:"id"`
	Code          string    `json:"code"`
	Description   *string   `json:"description"`
	DiscountType  string    `json:"discount_type"`
	DiscountValue float64   `json:"discount_value"`
	ValidFrom     time.Time `json:"valid_from"`
	ValidUntil    time.Time `json:"valid_until"`
	UsageLimit    *int      `json:"usage_limit"`
	UsedCount     int       `json:"used_count"`
}
