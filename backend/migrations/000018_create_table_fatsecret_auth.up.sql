CREATE TABLE FatsecretAuth (
    user_id BIGINT PRIMARY KEY REFERENCES Users (id),
    access_token TEXT NOT NULL,
    access_secret TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);