CREATE TABLE WorkoutExercises (
    id SERIAL PRIMARY KEY,
    workout_id BIGINT REFERENCES Workouts(id),
    exercise_id BIGINT REFERENCES Exercises(id),
    sets BIGINT DEFAULT NULL CHECK (sets > 0),
    reps BIGINT DEFAULT NULL CHECK (reps > 0),
    weight numeric(5,1) DEFAULT NULL CHECK (weight >= 0),
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);