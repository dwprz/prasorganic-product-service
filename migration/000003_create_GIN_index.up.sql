DROP INDEX IF EXISTS product_name_hash_index;

CREATE INDEX IF NOT EXISTS product_name_fts ON products USING GIN ((to_tsvector('indonesian', product_name)));