INSERT INTO users (name, email, password, role) VALUES ('usuario_anonimo', '', '', 'customer');
INSERT INTO categories(name, description) VALUES('Clásicos', 'Cuadros clásicos de diferentes estilos y tamaños');

-- // PRODUCT ITEMS
INSERT INTO products (category_id, slug, name, description, images, created_at, updated_at)
VALUES 
(1, 'la-noche-estrellada', 'La Noche Estrellada', 'Pintura famosa de Vincent van Gogh', '[{"url": "https://canvaslab.com/cdn/shop/files/473-Canvas_be7ffba8-d432-4f9c-a06a-42c47ab70528_500x.jpg?v=1720276978", "is_primary": true}]', DEFAULT, DEFAULT),
(1, 'el-grito', 'El Grito', 'Pintura icónica de Edvard Munch', '[{"url": "https://canvaslab.com/cdn/shop/products/100-Canvas_500x.jpg?v=1702506773", "is_primary": true}]', DEFAULT, DEFAULT);

INSERT INTO product_sizes (product_id, is_available, size_slug, size, price, discount, created_at, updated_at)
VALUES 
(1, TRUE, '20x30', '20 x 30', 100.00, NULL, DEFAULT, DEFAULT),
(1, TRUE, '50x100', '50 x 100', 200.00, 10.00, DEFAULT, DEFAULT),
(2, TRUE, '20x30', '20 x 30', 150.00, NULL, DEFAULT, DEFAULT),
(2, TRUE, '50x100', '50 x 100', 250.00, 15.00, DEFAULT, DEFAULT);

INSERT INTO products (category_id, slug, name, description, images, created_at, updated_at)
VALUES 
(1, 'la-ultima-cena', 'La Última Cena', 'Famosa pintura de Leonardo da Vinci', '[{"url": "https://canvaslab.com/cdn/shop/products/67523-Canvas_500x.jpg?v=1700268939", "is_primary": true}]', DEFAULT, DEFAULT),
(1, 'la-mona-lisa', 'La Mona Lisa', 'Retrato famoso de Leonardo da Vinci también conocido como La Mona Lisa', '[{"url": "https://canvaslab.com/cdn/shop/products/3-Canvas_651e9610-0105-4d39-bb99-22b3d6516306_500x.jpg?v=1703786723", "is_primary": true}]', DEFAULT, DEFAULT);

INSERT INTO product_sizes (product_id, is_available, size_slug, size, price, discount, created_at, updated_at)
VALUES 
(1, TRUE, '60x80', '60 x 80', 150.00, NULL, DEFAULT, DEFAULT),
(2, TRUE, '100x120', '50 x 100', 300.00, 199.00, DEFAULT, DEFAULT),
(3, TRUE, '20x30', '20 x 30', 150.00, NULL, DEFAULT, DEFAULT),
(3, TRUE, '50x100', '50 x 100', 250.00, 15.00, DEFAULT, DEFAULT),
(3, TRUE, '30x40', '30 x 40', 120.00, NULL, DEFAULT, DEFAULT),
(4, TRUE, '60x80', '60 x 80', 240.00, 20.00, DEFAULT, DEFAULT),
(4, TRUE, '25x35', '25 x 35', 130.00, NULL, DEFAULT, DEFAULT),
(4, TRUE, '70x100', '70 x 100', 270.00, 25.00, DEFAULT, DEFAULT);