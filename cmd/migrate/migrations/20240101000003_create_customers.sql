-- +goose Up
CREATE TABLE customers (
    id UUID PRIMARY KEY NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE customers;
