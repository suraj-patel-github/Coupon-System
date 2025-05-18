package coupon

import (
	"context"
	"encoding/json"
	"net/http"

	http2 "github.com/go-kit/kit/transport/http"
)

func NewHTTPHandler(endpoints Endpoints) http.Handler {
    r := http.NewServeMux()

    basePath := "/v1/coupons"

    r.Handle(basePath+"/admin", http2.NewServer(
        endpoints.CreateCouponEndpoint,
        decodeCreateCouponRequest,
        encodeResponse,
    ))

    r.Handle(basePath+"/applicable", http2.NewServer(
        endpoints.GetApplicableCouponsEndpoint,
        decodeApplicableCouponsRequest,
        encodeResponse,
    ))

    r.Handle(basePath+"/validate", http2.NewServer(
        endpoints.ValidateCouponEndpoint,
        decodeValidateCouponRequest,
        encodeResponse,
    ))

    return r
}

func decodeCreateCouponRequest(_ context.Context, r *http.Request) (interface{}, error) {
    var req Coupon
    err := json.NewDecoder(r.Body).Decode(&req)
    return req, err
}

func decodeApplicableCouponsRequest(_ context.Context, r *http.Request) (interface{}, error) {
    var req ApplicableRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    return req, err
}

func decodeValidateCouponRequest(_ context.Context, r *http.Request) (interface{}, error) {
    var req ValidateRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    return req, err
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
    w.Header().Set("Content-Type", "application/json")
    return json.NewEncoder(w).Encode(response)
}
