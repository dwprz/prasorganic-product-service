CREATE TABLE products (
    product_id SERIAL NOT NULL,
    product_name VARCHAR(100) NOT NULL,
    image_id VARCHAR(100) NOT NULL,
    image VARCHAR(500) NOT NULL,
    rating REAL,
    sold INTEGER,
    price INTEGER NOT NULL,
    stock INTEGER NOT NULL,
    length INTEGER NOT NULL,
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    weight REAL NOT NULL,
    description TEXT,
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP(3),
    CONSTRAINT products_pkey PRIMARY KEY (product_id), 
    CONSTRAINT product_name_unique UNIQUE (product_name)
);

CREATE INDEX IF NOT EXISTS product_name_hash_index ON products USING HASH (product_name);