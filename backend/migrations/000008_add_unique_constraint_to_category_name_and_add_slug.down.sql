ALTER TABLE Categories
    DROP CONSTRAINT unique_category_name,
    DROP COLUMN slug;