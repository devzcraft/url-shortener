CREATE TABLE IF NOT EXISTS urls(
    id INTEGER PRIMARY KEY,
    alias TEXT NOT NULL UNIQUE,
    url TEXT NOT NULL);
CREATE INDEX IF NOT EXISTS idx_alias ON urls(alias);