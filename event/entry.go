package event

import "time"

type Event struct {
	ID         string
	CalendarID string
	Status     string
	StartsAt   time.Time
	EndsAt     time.Time
}
