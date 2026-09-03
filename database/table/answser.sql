
CREATE TABLE IF NOT EXISTS questions_answers (
    id SERIAL PRIMARY KEY,
    profile_id INTEGER NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    option_id INTEGER NOT NULL REFERENCES answer_options(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW()
);