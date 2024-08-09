ALTER TABLE products ADD COLUMN category VARCHAR(20);

UPDATE products SET category = 'default_value' WHERE category IS NULL;

ALTER TABLE products ALTER COLUMN category SET NOT NULL;

CREATE INDEX IF NOT EXISTS category_hash_index ON products USING HASH (category);