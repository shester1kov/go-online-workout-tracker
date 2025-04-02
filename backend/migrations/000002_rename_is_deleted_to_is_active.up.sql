ALTER TABLE Exercises RENAME COLUMN is_deleted TO is_active;
ALTER TABLE Exercises ALTER COLUMN is_active SET DEFAULT TRUE;