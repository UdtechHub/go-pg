CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT UUID_GENERATE_V4() UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    price int,
    manufacturer_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMPTZ,
    FOREIGN KEY (manufacturer_id) REFERENCES manufacturers(id)
);