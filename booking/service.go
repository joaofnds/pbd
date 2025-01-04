package booking

import (
	"app/adapter/time"
	"app/calendar"
	"app/event"
	"app/internal/clock"
	"app/worker"
	"context"
)

type Service struct {
	clock     clock.Clock
	workers   *worker.Service
	calendars *calendar.Service
	events    *event.Service
}

func NewService(
	clock *time.Clock,
	workers *worker.Service,
	calendars *calendar.Service,
	events *event.Service,
) *Service {
	return &Service{
		clock:     clock,
		workers:   workers,
		calendars: calendars,
		events:    events,
	}
}

func (service *Service) Book(ctx context.Context, evt event.Event, request Request) (event.Event, error) {
	if err := service.ValidateRequest(evt, request); err != nil {
		return event.Event{}, err
	}

	deleteErr := service.events.Delete(ctx, evt)
	if deleteErr != nil {
		return event.Event{}, deleteErr
	}

	newEvent, createErr := service.events.Create(ctx, event.CreateDTO{
		CalendarID: evt.CalendarID,
		Status:     event.StatusBooked,
		StartsAt:   request.StartsAt,
		EndsAt:     request.EndsAt,
	})

	return newEvent, createErr
}

func (service *Service) ValidateRequest(evt event.Event, request Request) error {
	if request.StartsAt.IsZero() {
		return ErrMissingStartsAt
	}

	if request.EndsAt.IsZero() {
		return ErrMissingEndsAt
	}

	if request.StartsAt.After(request.EndsAt) {
		return ErrStartAfterEnd
	}

	if request.StartsAt.Equal(request.EndsAt) {
		return ErrStartEqualEnd
	}

	if request.StartsAt.Before(service.clock.Now()) {
		return ErrStartAfterNow
	}

	if evt.Status != event.StatusAvailable {
		return ErrEventNotAvailable
	}

	if request.StartsAt.Before(evt.StartsAt) || request.EndsAt.After(evt.EndsAt) {
		return ErrNotWithinEventTime
	}

	durationInHours := request.EndsAt.Sub(request.StartsAt).Hours()
	if durationInHours != float64(int(durationInHours)) || durationInHours < 1 || durationInHours > 8 {
		return ErrInvalidDuration
	}

	return nil
}
