CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
    original TEXT NOT NULL,
    shortcode VARCHAR(16) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    hits INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS access_logs (
    id SERIAL PRIMARY KEY,
    url_id INT REFERENCES urls(id),
    accessed_at TIMESTAMP NOT NULL,
    ip VARCHAR(64),
    action VARCHAR(16)
);