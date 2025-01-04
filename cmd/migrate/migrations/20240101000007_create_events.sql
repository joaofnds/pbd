-- +goose Up
CREATE TABLE events (
    id UUID PRIMARY KEY,
    calendar_id UUID NOT NULL REFERENCES calendars (id) ON DELETE CASCADE,
    status EVENT_STATUS NOT NULL,
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX ON events (status, starts_at, ends_at);

-- +goose Down
DROP TABLE events;
