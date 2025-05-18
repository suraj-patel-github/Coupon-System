package coupon

import (
	"coupon-system/types"
	"fmt"
)

type Coupon struct {
	ID                   int            `json:"id"` // Usually set by the DB, can be omitted in creation
	CouponCode           string         `json:"couponCode"`
	ExpiryDate           types.DateOnly `json:"expiryDate"`
	UsageType            string         `json:"usageType"`
	MinOrderValue        float64        `json:"minOrderValue"`
	ValidFrom            types.DateOnly `json:"validFrom"`
	ValidUntil           types.DateOnly `json:"validUntil"`
	DiscountType         string         `json:"discountType"`
	DiscountValue        float64        `json:"discountValue"`
	MaxUsagePerUser      int            `json:"maxUsagePerUser"`
	TermsAndConditions   string         `json:"termsAndConditions"`
	DiscountTarget       string         `json:"discountTarget"`
	ApplicableMedicines  []string       `json:"applicableMedicines"`
	ApplicableCategories []string       `json:"applicableCategories"`
}

type CartItem struct {
	Medicine string `json:"medicine"`
	Category string `json:"category"`
}

type ValidateRequest struct {
	CouponCode string         `json:"couponCode"`
	CartItems  []CartItem     `json:"cartItems"`
	OrderTotal float64        `json:"orderTotal"`
	Timestamp  types.DateOnly `json:"timestamp"`
}

type ValidateResponse struct {
	IsValid  bool   `json:"isValid"`
	Message  string `json:"message,omitempty"`
	Discount struct {
		ItemsDiscount   float64 `json:"itemsDiscount,omitempty"`
		ChargesDiscount float64 `json:"chargesDiscount,omitempty"`
	} `json:"discount,omitempty"`
}

type SuccessResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type CommonError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e CommonError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}
