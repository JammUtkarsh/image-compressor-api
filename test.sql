-- Inserting users
INSERT INTO Users (name, mobile, latitude, longitude)
VALUES ('John Doe', '1234567890', 40.7128, -74.0060),
       ('Jane Smith', '9876543210', 34.0522, -118.2437),
       ('Alice Johnson', '5551234567', 51.5074, -0.1278);

-- Inserting products for the users
INSERT INTO Products (user_id, product_name, product_description, product_images, product_price, compressed_product_images)
VALUES (1, 'Product A', 'Description for Product A', ARRAY['https://example.com/image1.jpg', 'https://example.com/image2.jpg'], 99.99, ARRAY['https://example.com/compressed_image1.jpg']),
       (1, 'Product B', 'Description for Product B', ARRAY['https://example.com/image3.jpg', 'https://example.com/image4.jpg'], 49.99, ARRAY['https://example.com/compressed_image2.jpg']),
       (2, 'Product C', 'Description for Product C', ARRAY['https://example.com/image5.jpg', 'https://example.com/image6.jpg'], 149.99, ARRAY['https://example.com/compressed_image3.jpg']),
       (3, 'Product D', 'Description for Product D', ARRAY['https://example.com/image7.jpg', 'https://example.com/image8.jpg'], 199.99, ARRAY['https://example.com/compressed_image4.jpg']);

-- Select all columns from Users table
SELECT * FROM Users WHERE id = 1;

-- Select all columns from Products table
SELECT * FROM Products WHERE product_id = 2;

-- To see triggers in action
SELECT pg_sleep(5);

-- Update mobile for a user
UPDATE Users SET mobile = '5556667777' WHERE id = 1;

-- Update product description for a product
UPDATE Products SET product_description = 'Updated description for Product B' WHERE product_id = 2;

-- Select all columns from Users table
SELECT * FROM Users WHERE id = 1;

-- Select all columns from Products table
SELECT * FROM Products WHERE product_id = 2;