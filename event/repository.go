package event

import (
	"context"
	"time"
)

type Repository interface {
	Create(context.Context, Event) error
	List(context.Context, ListDTO) ([]Event, error)
	FindByID(context.Context, string) (Event, error)
	Update(context.Context, Event, UpdateDTO) error
	Delete(context.Context, Event) error

	IsTimeSlotTaken(ctx context.Context, calendarID string, startsAt time.Time, endsAt time.Time) bool
}
