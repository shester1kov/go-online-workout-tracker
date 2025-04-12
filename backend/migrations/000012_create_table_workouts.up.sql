CREATE TABLE Workouts (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES Users(id),
    date DATE,
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    is_active BOOL DEFAULT TRUE
);