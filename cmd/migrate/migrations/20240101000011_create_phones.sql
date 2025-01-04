-- +goose Up
CREATE TABLE phones (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    country_code VARCHAR(3) NOT NULL,
    area_code VARCHAR(3) NOT NULL,
    number VARCHAR(12) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX ON phones (country_code, area_code, number);

-- +goose Down
DROP TABLE phones;
