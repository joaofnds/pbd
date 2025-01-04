package event

import "time"

type CreateDTO struct {
	CalendarID string    `json:"calendar_id"`
	Status     string    `json:"status"`
	StartsAt   time.Time `json:"starts_at"`
	EndsAt     time.Time `json:"ends_at"`
}

type UpdateDTO struct {
	Status string `json:"status"`
}

type ListDTO struct {
	CalendarID string
	StartsAt   time.Time
	EndsAt     time.Time
}
