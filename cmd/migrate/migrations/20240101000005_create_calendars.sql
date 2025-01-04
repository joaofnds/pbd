-- +goose Up
CREATE TABLE calendars (
    id UUID PRIMARY KEY REFERENCES workers (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE calendars;
