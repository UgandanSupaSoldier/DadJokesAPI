CREATE TABLE tokens (
    token TEXT NOT NULL,
    expires_at TIMESTAMP
);
CREATE UNIQUE INDEX token_index ON tokens(token);

CREATE TABLE jokes (
    id SERIAL PRIMARY KEY NOT NULL,
    text TEXT NOT NULL,
    author VARCHAR(255),
    category VARCHAR(255),
    rating NUMERIC(3,1),
    inserted_at TIMESTAMP NOT NULL DEFAULT now()
);
