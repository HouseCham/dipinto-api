create table users (
    id bigint primary key generated always as identity,
    name text not null,
    email text not null unique,
    phone text unique,
    password text not null,
    role text check (role in ('customer', 'admin')) not null,
    created_at timestamptz default now(),
    updated_at timestamptz default now(),
    deleted_at timestamptz
);

create table categories (
    id bigint primary key generated always as identity,
    name text not null,
    description text, 
    created_at timestamptz default now(),
    updated_at timestamptz default now()
);

create table products (
    id bigint primary key generated always as identity,
    category_id bigint references categories (id),
    slug text unique,
    name text not null,
    description text,
    images jsonb,
    created_at timestamptz default now(),
    updated_at timestamptz default now(),
    deleted_at timestamptz
);

create table product_sizes (
    id bigint primary key generated always as identity,
    product_id bigint references products (id),
    is_available boolean not null,
    size_slug text not null,
    size text not null,
    price numeric(10, 2) not null,
    discount numeric,
    created_at timestamptz default now(),
    updated_at timestamptz default now(),
    deleted_at timestamptz
);

CREATE TABLE coupons (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    code TEXT UNIQUE NOT NULL,
    description TEXT,
    discount_type TEXT CHECK (discount_type IN ('percentage', 'fixed')) NOT NULL,
    discount_value NUMERIC(10, 2) NOT NULL,
    valid_from DATE NOT NULL,
    valid_until DATE NOT NULL,
    usage_limit INT,
    used_count INT DEFAULT 0
);

create table addresses (
    id bigint primary key generated always as identity,
    user_id bigint references users (id),
    alias text,
    reference text,
    addressee text,
    phone text,
    addressee_email text,
    street_number text not null,
    department text,
    neighborhood text,
    city text not null,
    state text not null,
    postal_code text not null,
    country text not null,
    created_at timestamptz default now(),
    updated_at timestamptz default now(),
    deleted_at timestamptz
);

create table orders (
    id bigint primary key generated always as identity,
    user_id bigint references users (id) not null,
    address_id bigint references addresses (id) not null,
    coupon_id bigint references coupons(id),
    order_date timestamptz default now(),
    delivery_date timestamptz,
    total_amount numeric(10, 2) not null,
    tracking_id text,
    delivery_cost numeric(10, 2),
    shipping_company text,
    status text check (
        status in (
            'pending',
            'shipped',
            'delivered',
            'cancelled'
        )
    ) not null,
    payment_method text check (
        status in (
            'cash',
            'card'
        )
    ) not null
);

create table order_items (
    id bigint primary key generated always as identity,
    order_id bigint references orders (id),
    product_id bigint references product_sizes (id),
    quantity int not null,
    price numeric(10, 2) not null,
    discount numeric
);

create table reviews (
    id bigint primary key generated always as identity,
    product_id bigint references products (id),
    user_id bigint references users (id),
    rating int check (rating between 1 and 5) not null,
    comment text,
    created_at timestamptz default now()
);

create table payments (
    id bigint primary key generated always as identity,
    order_id bigint references orders (id),
    payment_method text not null,
    payment_date timestamptz default now(),
    amount numeric(10, 2) not null
);

create table expenses (
  id bigint primary key generated always as identity,
  amount numeric(10, 2) not null,
  expense_date date not null,
  description text,
  supplier text,
  ticker_url text
);

/* ========== CARTS  ========== */
create table carts (
  id bigint primary key generated always as identity,
  user_id bigint references users (id),
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);

create table cart_products (
  id bigint primary key generated always as identity,
  cart_id bigint references carts (id),
  product_id bigint references product_sizes (id),
  quantity int not null,
  added_at timestamptz default now()
);

/* ========== WISHLISTS  ========== */
create table wishlists (
  id bigint primary key generated always as identity,
  user_id bigint references users (id),
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);

create table wishlist_products (
  id bigint primary key generated always as identity,
  wishlist_id bigint references wishlists (id),
  product_id bigint references products (id),
  added_at timestamptz default now()
);

--Indexes

create index idx_products_name on products using btree (name);

create index idx_products_category_id on products using btree (category_id);

create index idx_product_sizes_product_id on product_sizes using btree (product_id);

create index idx_product_sizes_size on product_sizes using btree (size);

create index idx_users_email on users using btree (email);

create index idx_users_role on users using btree (role);

create index idx_expenses_expense_date on expenses (expense_date);

create index idx_expenses_supplier on expenses (supplier);

create index idx_carts_user_id on carts using btree (user_id);

create index idx_cart_products_cart_id on cart_products using btree (cart_id);

create index idx_wishlists_user_id on wishlists using btree (user_id);

create index idx_wishlist_products_wishlist_id on wishlist_products using btree (wishlist_id);

create table change_logs (
    id bigint primary key generated always as identity,
    table_name text not null,
    operation text check (
        operation in ('UPDATE', 'DELETE')
    ) not null,
    record_id bigint not null,
    changed_at timestamptz default now(),
    user_id bigint references users (id),
    details jsonb
);

/* Trigger function to log changes */
CREATE OR REPLACE FUNCTION log_changes() RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'UPDATE') THEN
        INSERT INTO change_logs (table_name, operation, record_id, user_id, details)
        VALUES (TG_TABLE_NAME, TG_OP, OLD.id, NULL, row_to_json(OLD));
        RETURN NEW;
    ELSIF (TG_OP = 'DELETE') THEN
        INSERT INTO change_logs (table_name, operation, record_id, user_id, details)
        VALUES (TG_TABLE_NAME, TG_OP, OLD.id, NULL, row_to_json(OLD));
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

/* Assign the trigger to all tables (delete to all, update to products) */
-- Products
CREATE TRIGGER log_product_updates
AFTER
UPDATE ON products FOR EACH ROW
EXECUTE FUNCTION log_changes ();

CREATE TRIGGER log_product_deletes
AFTER DELETE ON products FOR EACH ROW
EXECUTE FUNCTION log_changes ();

-- Product Sizes
CREATE TRIGGER log_product_size_deletes
AFTER DELETE ON product_sizes FOR EACH ROW
EXECUTE FUNCTION log_changes ();

-- Users
CREATE TRIGGER log_user_deletes
AFTER DELETE ON users FOR EACH ROW
EXECUTE FUNCTION log_changes ();

-- Categories
CREATE TRIGGER log_category_deletes
AFTER DELETE ON categories FOR EACH ROW
EXECUTE FUNCTION log_changes ();

-- Orders
CREATE TRIGGER log_order_deletes
AFTER DELETE ON orders FOR EACH ROW
EXECUTE FUNCTION log_changes ();

-- Order Items
CREATE TRIGGER log_order_item_deletes
AFTER DELETE ON order_items FOR EACH ROW
EXECUTE FUNCTION log_changes ();

-- Reviews
CREATE TRIGGER log_review_deletes
AFTER DELETE ON reviews FOR EACH ROW
EXECUTE FUNCTION log_changes ();

-- Addresses
CREATE TRIGGER log_address_deletes
AFTER DELETE ON addresses FOR EACH ROW
EXECUTE FUNCTION log_changes ();

-- Payments
CREATE TRIGGER log_payment_deletes
AFTER DELETE ON payments FOR EACH ROW
EXECUTE FUNCTION log_changes ();

-- Change Logs
CREATE TRIGGER log_change_log_deletes
AFTER DELETE ON change_logs FOR EACH ROW
EXECUTE FUNCTION log_changes ();

-- Expenses
CREATE TRIGGER log_expense_deletes
AFTER DELETE ON expenses FOR EACH ROW
EXECUTE FUNCTION log_changes ();

-- CREATE TRIGGER for carts, cart_products, wishlists, wishlist_products
CREATE TRIGGER log_cart_deletes
AFTER DELETE ON carts FOR EACH ROW
EXECUTE FUNCTION log_changes ();

CREATE TRIGGER log_wishlist_deletes
AFTER DELETE ON wishlists FOR EACH ROW
EXECUTE FUNCTION log_changes ();