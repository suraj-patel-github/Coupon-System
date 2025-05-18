package coupon

import (
	"context"
	"database/sql"
	"net/http"
	"sync"
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

// CreateCoupon inserts a new coupon into the database
func (s *couponService) CreateCoupon(ctx context.Context, c Coupon) error {
	return s.repo.CreateCoupon(c)
}

// GetApplicableCoupons filters coupons in parallel using goroutines and channels
func (s *couponService) GetApplicableCoupons(ctx context.Context, cartItems []CartItem, orderTotal float64, now time.Time) ([]Coupon, error) {
	coupons, err := s.repo.GetAllCoupons()
	if err != nil {
		return nil, CommonError{Code: http.StatusInternalServerError, Message: "Internal Server Error"}
	}

	var applicable []Coupon
	ch := make(chan Coupon)
	var wg sync.WaitGroup

	for _, coupon := range coupons {
		wg.Add(1)
		go func(c Coupon) {
			defer wg.Done()
			if validateCoupon(c, cartItems, orderTotal, now) {
				ch <- c
			}
		}(coupon)
	}

	// Close channel after all goroutines finish
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Collect from channel
	for c := range ch {
		applicable = append(applicable, c)
	}

	return applicable, nil
}

// ValidateCoupon checks if a single coupon is valid for the given request
func (s *couponService) ValidateCoupon(ctx context.Context, req ValidateRequest) (ValidateResponse, error) {
	coupon, err := s.repo.GetCouponByCode(req.CouponCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return ValidateResponse{
				IsValid: false,
				Message: "Coupon not found",
			}, nil
		}
		return ValidateResponse{
			IsValid: false,
			Message: "Internal server error",
		}, CommonError{Code: http.StatusInternalServerError, Message: "internal server error"}

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
