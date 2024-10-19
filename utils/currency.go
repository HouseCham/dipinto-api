package utils

import "github.com/HouseCham/dipinto-api/internal/domain/model"

// ApplyCouponDiscount applies the discount value of the coupon to the total amount
func ApplyCouponDiscount(totalAmount float64, quantity float64, coupon *model.Coupon) float64 {
	if coupon.DiscountType == "percentage" {
		return totalAmount - (totalAmount * coupon.DiscountValue / 100)
	}
	return totalAmount - coupon.DiscountValue
}