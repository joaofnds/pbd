package worker

import (
	"app/adapter/time"
	"app/internal/clock"
	"context"
)

type Service struct {
	clock       clock.Clock
	repo        Repository
	permissions *PermissionService
}

func NewService(
	clock *time.Clock,
	repo Repository,
	permissions *PermissionService,
) *Service {
	return &Service{
		clock:       clock,
		repo:        repo,
		permissions: permissions,
	}
}

func (service *Service) Create(ctx context.Context, userID string, hourlyRate int) (Worker, error) {
	worker := Worker{
		ID:         userID,
		HourlyRate: hourlyRate,
		CreatedAt:  service.clock.Now(),
		UpdatedAt:  service.clock.Now(),
	}

	if err := service.repo.Create(ctx, worker); err != nil {
		return Worker{}, err
	}

	return worker, service.permissions.GrantNewWorkerPermissions(worker)
}

func (service *Service) FindByID(ctx context.Context, id string) (Worker, error) {
	return service.repo.FindByID(ctx, id)
}
