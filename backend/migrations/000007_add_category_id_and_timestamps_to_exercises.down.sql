ALTER TABLE Exercises
    DROP COLUMN category_id,

    ALTER COLUMN name TYPE TEXT,

    DROP COLUMN created_at,
    DROP COLUMN updated_at,

    ADD COLUMN category TEXT NOT NULL;