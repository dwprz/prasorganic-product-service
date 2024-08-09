ALTER TABLE products DROP COLUMN category;

DROP INDEX IF EXISTS category_hash_index;