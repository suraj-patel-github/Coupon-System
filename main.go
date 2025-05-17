package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"coupon-system/coupon"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:password@localhost:5432/coupons?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	repository := coupon.NewPostgresRepository(db)
	service := coupon.NewService(repository)
	endpoints := coupon.MakeEndpoints(service)
	handler := coupon.NewHTTPHandler(endpoints)

	r := mux.NewRouter()
	r.PathPrefix("/").Handler(handler)

	fmt.Println("Coupon service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
