
CREATE TABLE IF NOT EXISTS profiles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    date_of_birth DATE NOT NULL,
    skills TEXT[],
    sector VARCHAR(100),
    location VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW()
);
