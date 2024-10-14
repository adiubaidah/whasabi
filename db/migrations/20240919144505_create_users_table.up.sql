CREATE TYPE ROLE_USER AS ENUM ('admin', 'user');
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    role ROLE_USER NOT NULL,
    password VARCHAR(255) NOT NULL,
    is_active boolean NOT NULL DEFAULT FALSE
);