CREATE TABLE questions_answers (
    id SERIAL PRIMARY KEY,
    profile_id INTEGER NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    options TEXT[],
    created_at TIMESTAMP DEFAULT NOW()
);
