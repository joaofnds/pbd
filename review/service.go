package review

import (
	"app/internal/clock"
	"app/internal/id"
	"app/order"
	"context"
	"errors"
)

type Service struct {
	id    id.Generator
	clock clock.Clock
	repo  Repository
}

func NewService(id id.Generator, clock clock.Clock, repo Repository) *Service {
	return &Service{
		id:    id,
		clock: clock,
		repo:  repo,
	}
}

func (service *Service) Create(ctx context.Context, ord order.Order, rating int, comment string) (Review, error) {
	if ord.Status != order.StatusCompleted {
		return Review{}, ErrOrderNotCompleted
	}

	_, getErr := service.GetByOrderID(ctx, ord.ID)
	switch {
	case errors.Is(getErr, ErrNotFound):
		// proceed
	case getErr == nil:
		return Review{}, ErrOrderAlreadyReviewed
	default:
		return Review{}, getErr
	}

	review := Review{
		ID:        service.id.NewID(),
		OrderID:   ord.ID,
		Rating:    rating,
		Comment:   comment,
		CreatedAt: service.clock.Now(),
		UpdatedAt: service.clock.Now(),
	}

	return review, service.repo.Create(ctx, review)
}

func (service *Service) GetByOrderID(ctx context.Context, orderID string) (Review, error) {
	return service.repo.GetByOrderID(ctx, orderID)
}

func (service *Service) FindByWorkerID(ctx context.Context, workerID string) ([]Review, error) {
	return service.repo.FindByWorkerID(ctx, workerID)
}

func (service *Service) FindByCustomerID(ctx context.Context, orderID string) ([]Review, error) {
	return service.repo.FindByCustomerID(ctx, orderID)
}

func (service *Service) Delete(ctx context.Context, review Review) error {
	return service.repo.Delete(ctx, review)
}
