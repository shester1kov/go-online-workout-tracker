ALTER TABLE Categories
    ADD COLUMN slug VARCHAR(100) NOT NULL UNIQUE,
    ADD CONSTRAINT unique_category_name UNIQUE (name);