-- +goose Up
CREATE TABLE workers (
    id UUID PRIMARY KEY NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    hourly_rate INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE workers;
