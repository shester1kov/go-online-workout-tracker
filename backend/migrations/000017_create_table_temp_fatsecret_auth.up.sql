CREATE TABLE TempFatsecretAuth (
    user_id BIGINT PRIMARY KEY REFERENCES Users (id),
    request_token TEXT NOT NULL,
    request_secret TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);