--- USERS
CREATE TABLE Users (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255),
	mobile VARCHAR(20),
	latitude FLOAT,
	longitude FLOAT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_users_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_updated_at_trigger
BEFORE UPDATE ON Users
FOR EACH ROW
EXECUTE FUNCTION update_users_updated_at();

--- PRODUCTS
CREATE TABLE Products (
	product_id SERIAL PRIMARY KEY,
	user_id INT REFERENCES Users(id),
	product_name VARCHAR(255),
	product_description TEXT,
	product_images TEXT [],
	product_price NUMERIC,
	compressed_product_images TEXT [],
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_products_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER products_updated_at_trigger
BEFORE UPDATE ON Products
FOR EACH ROW
EXECUTE FUNCTION update_products_updated_at();


-- Insert a user
INSERT INTO Users (name, mobile, latitude, longitude) VALUES ('Admin', '1234567890', 12.9716, 77.5946);