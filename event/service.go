package event

import (
	"app/adapter/uuid"
	"app/internal/id"
	"context"
)

type Service struct {
	id   id.Generator
	repo Repository
}

func NewService(
	id *uuid.Generator,
	repo Repository,
) *Service {
	return &Service{
		id:   id,
		repo: repo,
	}
}

func (service *Service) Create(ctx context.Context, dto CreateDTO) (Event, error) {
	if service.repo.IsTimeSlotTaken(ctx, dto.CalendarID, dto.StartsAt, dto.EndsAt) {
		return Event{}, ErrTimeSlotTaken
	}

	evt := Event{
		ID:         service.id.NewID(),
		CalendarID: dto.CalendarID,
		Status:     dto.Status,
		StartsAt:   dto.StartsAt,
		EndsAt:     dto.EndsAt,
	}

	return evt, service.repo.Create(ctx, evt)
}

func (service *Service) List(ctx context.Context, dto ListDTO) ([]Event, error) {
	return service.repo.List(ctx, dto)
}

func (service *Service) FindByID(ctx context.Context, eventID string) (Event, error) {
	return service.repo.FindByID(ctx, eventID)
}

func (service *Service) Update(ctx context.Context, event Event, dto UpdateDTO) error {
	return service.repo.Update(ctx, event, dto)
}

func (service *Service) Delete(ctx context.Context, event Event) error {
	return service.repo.Delete(ctx, event)
}
