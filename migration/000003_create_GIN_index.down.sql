DROP INDEX IF EXISTS product_name_fts;

CREATE INDEX IF NOT EXISTS product_name_hash_index ON products USING HASH (product_name);