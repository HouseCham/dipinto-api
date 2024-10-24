-- Drop triggers for logging changes
DROP TRIGGER IF EXISTS log_product_updates ON products;
DROP TRIGGER IF EXISTS log_product_deletes ON products;
DROP TRIGGER IF EXISTS log_product_size_deletes ON product_sizes;
DROP TRIGGER IF EXISTS log_user_deletes ON users;
DROP TRIGGER IF EXISTS log_category_deletes ON categories;
DROP TRIGGER IF EXISTS log_order_deletes ON orders;
DROP TRIGGER IF EXISTS log_order_item_deletes ON order_items;
DROP TRIGGER IF EXISTS log_review_deletes ON reviews;
DROP TRIGGER IF EXISTS log_address_deletes ON addresses;
DROP TRIGGER IF EXISTS log_payment_deletes ON payments;
DROP TRIGGER IF EXISTS log_change_log_deletes ON change_logs;
DROP TRIGGER IF EXISTS log_expense_deletes ON expenses;
DROP TRIGGER IF EXISTS log_cart_deletes ON carts;
DROP TRIGGER IF EXISTS log_wishlist_deletes ON wishlists;


-- Drop dependent tables first
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS addresses;
DROP TABLE IF EXISTS cart_products;
DROP TABLE IF EXISTS wishlist_products;
DROP TABLE IF EXISTS product_sizes;

-- Drop primary tables
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS carts;
DROP TABLE IF EXISTS wishlists;
DROP TABLE IF EXISTS coupons;
DROP TABLE IF EXISTS expenses;
DROP TABLE IF EXISTS change_logs;
DROP TABLE IF EXISTS users;

-- Drop trigger function
DROP FUNCTION IF EXISTS log_changes;
