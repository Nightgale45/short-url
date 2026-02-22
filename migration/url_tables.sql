CREATE TABLE IF NOT EXISTS url_data (
    id BIGSERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    counter INT NOT NULL,
    passcode TEXT,
    salt BIGINT NOT NULL
);
