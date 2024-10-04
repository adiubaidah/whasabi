CREATE TYPE ROLE_AS AS ENUM ('model', 'user');
CREATE TABLE IF NOT EXISTS histories (
    id SERIAL PRIMARY KEY,
    process_id INT NOT NULL,
    sender VARCHAR(255) NOT NULL,
    receiver VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    role_as ROLE_AS NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (process_id) REFERENCES process (id)  ON DELETE CASCADE
);