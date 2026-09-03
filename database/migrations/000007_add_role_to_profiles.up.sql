CREATE TYPE user_role AS ENUM ('candidate', 'recruiter', 'admin');

ALTER TABLE profiles
    ADD COLUMN role user_role NOT NULL DEFAULT 'candidate';