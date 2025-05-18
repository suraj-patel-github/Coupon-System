package coupon

import (
	"time"
)

func validateCoupon(c Coupon, cartItems []CartItem, orderTotal float64, now time.Time) bool {
	// Check expiry

	if now.After(c.ExpiryDate.ToTime()) {
		return false
	}

	if !c.ValidFrom.ToTime().IsZero() && now.Before(c.ValidFrom.ToTime()) {
		return false
	}

	if !c.ValidUntil.ToTime().IsZero() && now.After(c.ValidUntil.ToTime()) {
		return false
	}

	// Check minimum order
	if c.MinOrderValue > 0 && orderTotal < c.MinOrderValue {
		return false
	}

	// Check applicable medicines or categories
	if len(c.ApplicableMedicines) > 0 {
		found := false
		for _, item := range cartItems {
			for _, med := range c.ApplicableMedicines {
				if item.Medicine == med {
					found = true
					break
				}
			}
		}
		if !found {
			return false
		}
	}

	if len(c.ApplicableCategories) > 0 {
		found := false
		for _, item := range cartItems {
			for _, cat := range c.ApplicableCategories {
				if item.Category == cat {
					found = true
					break
				}
			}
		}
		if !found {
			return false
		}
	}

	return true
}
