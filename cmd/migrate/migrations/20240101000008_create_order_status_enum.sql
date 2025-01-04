-- +goose Up
CREATE TYPE order_status AS ENUM (
    'created',
    'booked',
    'payment_pending',
    'payment_failed',
    'payment_succeeded',
    'canceled',
    'completed'
);

-- +goose Down
DROP TYPE order_status;
