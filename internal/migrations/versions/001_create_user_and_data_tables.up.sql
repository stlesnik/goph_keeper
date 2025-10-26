CREATE TABLE IF NOT EXISTS users
(
    id            VARCHAR PRIMARY KEY,
    username      VARCHAR NOT NULL,
    email         VARCHAR NOT NULL UNIQUE,
    password_hash VARCHAR NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE,
    updated_at    TIMESTAMP WITH TIME ZONE
);

CREATE TYPE type_enum AS ENUM ('password', 'text', 'binary', 'card');
CREATE TABLE IF NOT EXISTS data
(
    id             VARCHAR PRIMARY KEY,
    user_id        VARCHAR   NOT NULL,
    type           type_enum NOT NULL,
    title          VARCHAR   NOT NULL,
    encrypted_data VARCHAR   NOT NULL,
    iv             VARCHAR   NOT NULL,
    metadata       VARCHAR   NOT NULL,
    created_at     TIMESTAMP WITH TIME ZONE,
    updated_at     TIMESTAMP WITH TIME ZONE,

    CONSTRAINT users FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
);