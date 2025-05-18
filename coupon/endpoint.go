package coupon

import (
	"context"
	"coupon-system/utils"
	"net/http"

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
            er := err.(CommonError)
            return nil, er
        }
        return SuccessResponse{Code: http.StatusCreated, Message: "coupon created Successfully"}, nil
    }
}

type ApplicableRequest struct {
    CartItems  []CartItem `json:"cartItems"`
    OrderTotal float64    `json:"orderTotal"`
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
                "couponCode":    c.CouponCode,
                "discountValue": c.DiscountValue,
            })
        }

        return map[string]interface{}{"applicableCoupons": response}, nil
    }
}

func makeValidateCouponEndpoint(svc Service) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        req := request.(ValidateRequest)
        return svc.ValidateCoupon(ctx, req)
    }
}
