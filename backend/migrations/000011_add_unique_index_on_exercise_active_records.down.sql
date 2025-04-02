DROP INDEX unique_exercise_name;
ALTER TABLE exercises ADD CONSTRAINT unique_exercise_name UNIQUE (name);