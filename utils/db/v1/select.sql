-- Active: 1728327431411@@127.0.0.1@5432@dipinto_db@public
SELECT cp.id, p.name, p.images, p.slug, ps.size, ps.size_slug, ps.price, ps.discount, cp.quantity FROM cart_products cp
INNER JOIN product_sizes ps ON cp.product_id = ps.id
INNER JOIN products p ON ps.product_id = p.id
WHERE cp.cart_id = 1;

SELECT * FROM cart_products;

SELECT * FROM wishlists;
SELECT * FROM wishlist_products;

SELECT ps.id, p.name, p.slug, p.images, ps.size, ps.price, ps.discount FROM product_sizes ps
INNER JOIN products p ON ps.product_id = p.id
WHERE ps.product_id IN (1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
AND p.deleted_at IS NULL
AND ps.deleted_at IS NULL;

SELECT * FROM coupons;

SELECT * FROM product_sizes;

SELECT ps.id, p.name, p.slug, p.images, ps.size, ps.price, ps.discount 
FROM product_sizes ps 
INNER JOIN products p ON ps.product_id = p.id 
WHERE ps.id IN (20,22) 
AND ps.is_available = true 
AND ps.deleted_at IS NULL 
AND p.deleted_at IS NULL;