CREATE TABLE Exercises (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    description TEXT, 
    is_deleted BOOLEAN DEFAULT FALSE
);