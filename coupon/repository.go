package coupon

import (
	"database/sql"
	"net/http"

	"github.com/lib/pq"
)

type Repository interface {
	CreateCoupon(c Coupon) error
	GetAllCoupons() ([]Coupon, error)
	GetCouponByCode(code string) (Coupon, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db}
}

func (r *postgresRepository) CreateCoupon(c Coupon) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	var id int
	err = tx.QueryRow(`
        INSERT INTO coupons (coupon_code, expiry_date, usage_type, min_order_value,
        valid_from, valid_until, discount_type, discount_value,
        max_usage_per_user, terms_and_conditions, discount_target)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id
    `,
		c.CouponCode, c.ExpiryDate.ToTime(), c.UsageType, c.MinOrderValue,
		c.ValidFrom.ToTime(), c.ValidUntil.ToTime(), c.DiscountType, c.DiscountValue,
		c.MaxUsagePerUser, c.TermsAndConditions, c.DiscountTarget,
	).Scan(&id)
	if err != nil {
		tx.Rollback()
		// Check for duplicate coupon code
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" && pqErr.Constraint == "coupons_coupon_code_key" {
				return CommonError{Code: http.StatusBadRequest, Message: "Coupon already"}
			}
		}
		return err
	}

	for _, med := range c.ApplicableMedicines {
		_, err := tx.Exec(`INSERT INTO coupon_applicable_medicines (coupon_id, medicine_id) VALUES ($1, $2)`, id, med)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, cat := range c.ApplicableCategories {
		_, err := tx.Exec(`INSERT INTO coupon_applicable_categories (coupon_id, category) VALUES ($1, $2)`, id, cat)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *postgresRepository) GetAllCoupons() ([]Coupon, error) {
	rows, err := r.db.Query(`
    SELECT
      c.id, c.coupon_code, c.expiry_date, c.usage_type, c.min_order_value,
      c.valid_from, c.valid_until, c.discount_type, c.discount_value,
      c.max_usage_per_user, c.terms_and_conditions, c.discount_target,
      COALESCE(array_agg(DISTINCT cm.medicine_id), '{}') AS applicable_medicines,
      COALESCE(array_agg(DISTINCT cc.category), '{}') AS applicable_categories
    FROM coupons c
    LEFT JOIN coupon_applicable_medicines cm ON cm.coupon_id = c.id
    LEFT JOIN coupon_applicable_categories cc ON cc.coupon_id = c.id
    GROUP BY c.id
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coupons []Coupon
	for rows.Next() {
		var c Coupon
		err := rows.Scan(
			&c.ID, &c.CouponCode, &c.ExpiryDate, &c.UsageType,
			&c.MinOrderValue, &c.ValidFrom, &c.ValidUntil,
			&c.DiscountType, &c.DiscountValue, &c.MaxUsagePerUser,
			&c.TermsAndConditions, &c.DiscountTarget,
			pq.Array(&c.ApplicableMedicines),
			pq.Array(&c.ApplicableCategories),
		)
		if err != nil {
			return nil, err
		}
		coupons = append(coupons, c)
	}
	return coupons, nil

}

func (r *postgresRepository) GetCouponByCode(code string) (Coupon, error) {
	var c Coupon
	err := r.db.QueryRow(`
        SELECT id, coupon_code, expiry_date, usage_type, min_order_value,
        valid_from, valid_until, discount_type, discount_value, max_usage_per_user,
        terms_and_conditions, discount_target
        FROM coupons WHERE coupon_code = $1
    `, code).Scan(
		&c.ID, &c.CouponCode, &c.ExpiryDate, &c.UsageType, &c.MinOrderValue,
		&c.ValidFrom, &c.ValidUntil, &c.DiscountType, &c.DiscountValue, &c.MaxUsagePerUser,
		&c.TermsAndConditions, &c.DiscountTarget,
	)
	return c, err
}
