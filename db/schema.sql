CREATE TABLE coupons (
  id SERIAL PRIMARY KEY,
  coupon_code TEXT UNIQUE NOT NULL,
  expiry_date TIMESTAMP NOT NULL,
  usage_type TEXT NOT NULL CHECK (usage_type IN ('one_time', 'multi_use', 'time_based')),
  min_order_value NUMERIC,
  valid_from TIMESTAMP,
  valid_until TIMESTAMP,
  discount_type TEXT NOT NULL CHECK (discount_type IN ('fixed', 'percentage')),
  discount_value NUMERIC,
  max_usage_per_user INT,
  terms_and_conditions TEXT,
  discount_target TEXT NOT NULL CHECK (discount_target IN ('inventory', 'charges')),
  created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE coupon_applicable_medicines (
  coupon_id INT REFERENCES coupons(id) ON DELETE CASCADE,
  medicine_id TEXT NOT NULL
);

CREATE TABLE coupon_applicable_categories (
  coupon_id INT REFERENCES coupons(id) ON DELETE CASCADE,
  category TEXT NOT NULL
);

CREATE TABLE coupon_usage (
  id SERIAL PRIMARY KEY,
  coupon_id INT REFERENCES coupons(id) ON DELETE CASCADE,
  user_id TEXT NOT NULL,
  used_at TIMESTAMP DEFAULT now()
);
