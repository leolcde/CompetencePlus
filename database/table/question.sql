
CREATE TABLE IF NOT EXISTS questions (
    id SERIAL PRIMARY KEY,
    content VARCHAR(255) NOT NULL,
    options TEXT[],
    weight INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
)
