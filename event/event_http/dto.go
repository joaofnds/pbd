package event_http

import "time"

type CreateBody struct {
	Status   string    `json:"status" validate:"required,oneof=available booked canceled"`
	StartsAt time.Time `json:"starts_at" validate:"required"`
	EndsAt   time.Time `json:"ends_at" validate:"required"`
}

type UpdateBody struct {
	Status string `json:"status" validate:"required,oneof=available booked"`
}

type TimeQuery struct {
	StartsAt time.Time `query:"starts_at" validate:"required"`
	EndsAt   time.Time `query:"ends_at" validate:"required,gtfield=StartsAt"`
}
