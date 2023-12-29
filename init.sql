CREATE TABLE Users (
  id bigserial PRIMARY KEY,
  name varchar NOT NULL,
  mobile varchar,
  latitude float,
  longitude float,
  created_at timestamptz NOT NULL DEFAULT (now()),
  updated_at timestamptz NOT NULL DEFAULT (now())
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

CREATE TABLE Products (
  id bigserial PRIMARY KEY,
  user_id bigint REFERENCES Users(id),
  product_name varchar NOT NULL,
  product_description text NOT NULL,
  product_price real NOT NULL,
  product_images text[],
  compressed_product_images text[],
  created_at timestamptz NOT NULL DEFAULT (now()),
  updated_at timestamptz NOT NULL DEFAULT (now())
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

-- Indexes were pushed at the bottom of the file to avoid: ERROR: relation "Users" does not exist
CREATE INDEX ON Users (id);

CREATE INDEX ON Products (id);

CREATE INDEX ON Products (user_id);

CREATE INDEX ON Products (user_id, id);

-- Insert a user
INSERT INTO
  Users (name, mobile, latitude, longitude)
VALUES
  ('admin', '1234567890', 0, 0);