CREATE TABLE TempFatsecretAuth (
    user_id TEXT PRIMARY KEY,
    request_token TEXT NOT NULL,
    request_secret TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);