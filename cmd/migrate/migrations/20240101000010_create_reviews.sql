-- +goose Up
CREATE TABLE reviews (
    id UUID PRIMARY KEY,
    order_id UUID UNIQUE NOT NULL UNIQUE REFERENCES orders (
        id
    ) ON DELETE CASCADE,
    rating INT NOT NULL,
    comment TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE reviews;
