ALTER TABLE Exercises RENAME COLUMN is_active TO is_deleted;
ALTER TABLE Exercises ALTER COLUMN is_deleted SET DEFAULT TRUE;