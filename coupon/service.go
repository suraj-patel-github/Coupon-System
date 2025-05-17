package coupon

import (
    "context"
    "time"
)

type Service interface {
    CreateCoupon(ctx context.Context, c Coupon) error
    GetApplicableCoupons(ctx context.Context, cartItems []CartItem, orderTotal float64, now time.Time) ([]Coupon, error)
    ValidateCoupon(ctx context.Context, req ValidateRequest) (ValidateResponse, error)
}

type couponService struct {
    repo Repository
}

func NewService(r Repository) Service {
    return &couponService{repo: r}
}

func (s *couponService) CreateCoupon(ctx context.Context, c Coupon) error {
    return s.repo.CreateCoupon(c)
}

func (s *couponService) GetApplicableCoupons(ctx context.Context, cartItems []CartItem, orderTotal float64, now time.Time) ([]Coupon, error) {
    coupons, err := s.repo.GetAllCoupons()
    if err != nil {
        return nil, err
    }

    var applicable []Coupon
    for _, coupon := range coupons {
        valid := validateCoupon(coupon, cartItems, orderTotal, now)
        if valid {
            applicable = append(applicable, coupon)
        }
    }

    return applicable, nil
}

func (s *couponService) ValidateCoupon(ctx context.Context, req ValidateRequest) (ValidateResponse, error) {
    coupon, err := s.repo.GetCouponByCode(req.CouponCode)
    if err != nil {
        return ValidateResponse{
            IsValid: false,
            Message: "Coupon not found",
        }, err
    }

    if !validateCoupon(coupon, req.CartItems, req.OrderTotal, req.Timestamp) {
        return ValidateResponse{
            IsValid: false,
            Message: "Coupon is not applicable or expired",
        }, nil
    }

    var discountAmount float64
    if coupon.DiscountType == "fixed" {
        discountAmount = coupon.DiscountValue
    } else if coupon.DiscountType == "percentage" {
        discountAmount = (req.OrderTotal * coupon.DiscountValue) / 100
    }

    resp := ValidateResponse{
        IsValid: true,
        Message: "coupon applied successfully",
    }

    if coupon.DiscountTarget == "inventory" {
        resp.Discount.ItemsDiscount = discountAmount
    } else {
        resp.Discount.ChargesDiscount = discountAmount
    }

    return resp, nil
}
