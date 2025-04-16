CREATE TABLE Foods (
    id SERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES Users(id),
    date TIMESTAMP DEFAULT NOW(),
    name VARCHAR(100) NOT NULL,
    quantity FLOAT NOT NULL, 
    unit VARCHAR(20) NOT NULL, 
    weight_grams FLOAT NOT NULL, 
    calories FLOAT NOT NULL, 
    protein FLOAT NOT NULL, 
    carbs FLOAT NOT NULL, 
    fat FLOAT NOT NULL
);