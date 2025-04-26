CREATE TABLE FatsecretAuth (
    user_id TEXT PRIMARY KEY,
    access_token TEXT NOT NULL,
    access_secret TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);