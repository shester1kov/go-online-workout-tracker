ALTER TABLE categories DROP CONSTRAINT unique_category_name;
ALTER TABLE categories DROP CONSTRAINT categories_slug_key;
CREATE UNIQUE INDEX unique_category_name ON categories (name) WHERE is_active = true;
CREATE UNIQUE INDEX unique_category_slug ON categories (slug) WHERE is_active = true;