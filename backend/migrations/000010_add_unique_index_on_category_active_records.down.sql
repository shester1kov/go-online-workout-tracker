
DROP INDEX unique_category_name;
DROP INDEX unique_category_slug;
ALTER TABLE categories ADD CONSTRAINT unique_category_name UNIQUE (name);
ALTER TABLE categories ADD CONSTRAINT categories_slug_key UNIQUE (slug);