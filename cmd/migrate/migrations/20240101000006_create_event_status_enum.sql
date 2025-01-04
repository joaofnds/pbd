-- +goose Up
CREATE TYPE event_status AS ENUM (
    'available',
    'booked',
    'canceled'
);

-- +goose Down
DROP TYPE event_status;
