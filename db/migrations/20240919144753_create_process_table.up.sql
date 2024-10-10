CREATE TABLE IF NOT EXISTS process (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(13) NOT NULL UNIQUE,
    instruction TEXT NOT NULL,
    temperature FLOAT NOT NULL,
    top_k FLOAT DEFAULT 64,
    top_p FLOAT DEFAULT 0.95,
    is_authenticated BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT phone_format CHECK (phone ~ '^62[0-9]{9,11}$'),
    
    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- Create a function to update the updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create a trigger to call the function before any update
CREATE TRIGGER update_updated_at
BEFORE UPDATE ON process
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();