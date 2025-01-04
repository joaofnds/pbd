-- +goose Up
CREATE TABLE orders (
    id UUID PRIMARY KEY,

    price INT NOT NULL,
    status ORDER_STATUS NOT NULL,

    address_id UUID NOT NULL REFERENCES addresses (id) ON DELETE RESTRICT,
    event_id UUID UNIQUE NOT NULL REFERENCES events (id) ON DELETE RESTRICT,
    worker_id UUID NOT NULL REFERENCES workers (id) ON DELETE CASCADE,
    customer_id UUID NOT NULL REFERENCES customers (id) ON DELETE CASCADE,

    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE orders;
