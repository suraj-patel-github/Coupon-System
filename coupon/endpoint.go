package coupon

import (
	"context"
	"coupon-system/utils"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
    CreateCouponEndpoint        endpoint.Endpoint
    GetApplicableCouponsEndpoint endpoint.Endpoint
    ValidateCouponEndpoint      endpoint.Endpoint
}

func MakeEndpoints(svc Service) Endpoints {
    return Endpoints{
        CreateCouponEndpoint: makeCreateCouponEndpoint(svc),
        GetApplicableCouponsEndpoint: makeGetApplicableCouponsEndpoint(svc),
        ValidateCouponEndpoint: makeValidateCouponEndpoint(svc),
    }
}

func makeCreateCouponEndpoint(svc Service) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        req := request.(Coupon)
        err := svc.CreateCoupon(ctx, req)
        if err != nil {
            return map[string]string{"error": err.Error()}, nil
        }
        return map[string]string{"status": "created"}, nil
    }
}

type ApplicableRequest struct {
    CartItems  []CartItem `json:"cart_items"`
    OrderTotal float64    `json:"order_total"`
    Timestamp  string     `json:"timestamp"`
}

func makeGetApplicableCouponsEndpoint(svc Service) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        req := request.(ApplicableRequest)
        parsedTime, _ := utils.ParseTimestamp(req.Timestamp)
        coupons, err := svc.GetApplicableCoupons(ctx, req.CartItems, req.OrderTotal, parsedTime)
        if err != nil {
            return nil, err
        }

        var response []map[string]interface{}
        for _, c := range coupons {
            response = append(response, map[string]interface{}{
                "coupon_code":    c.CouponCode,
                "discount_value": c.DiscountValue,
            })
        }

        return map[string]interface{}{"applicable_coupons": response}, nil
    }
}

func makeValidateCouponEndpoint(svc Service) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        req := request.(ValidateRequest)
        return svc.ValidateCoupon(ctx, req)
    }
}
