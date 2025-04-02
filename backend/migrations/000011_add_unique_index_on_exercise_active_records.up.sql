ALTER TABLE exercises DROP CONSTRAINT unique_exercise_name;
CREATE UNIQUE INDEX unique_exercise_name ON exercises (name) WHERE is_active = true;