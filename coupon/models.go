package coupon

import "time"

type Coupon struct {
    ID                   int
    CouponCode           string
    ExpiryDate           time.Time
    UsageType            string
    MinOrderValue        float64
    ValidFrom            time.Time
    ValidUntil           time.Time
    DiscountType         string
    DiscountValue        float64
    MaxUsagePerUser      int
    TermsAndConditions   string
    DiscountTarget       string
    ApplicableMedicines  []string
    ApplicableCategories []string
}

type CartItem struct {
    ID       string `json:"id"`
    Category string `json:"category"`
}

type ValidateRequest struct {
    CouponCode string     `json:"coupon_code"`
    CartItems  []CartItem `json:"cart_items"`
    OrderTotal float64    `json:"order_total"`
    Timestamp  time.Time  `json:"timestamp"`
}

type ValidateResponse struct {
    IsValid bool   `json:"is_valid"`
    Message string `json:"message,omitempty"`
    Discount struct {
        ItemsDiscount   float64 `json:"items_discount,omitempty"`
        ChargesDiscount float64 `json:"charges_discount,omitempty"`
    } `json:"discount,omitempty"`
}
