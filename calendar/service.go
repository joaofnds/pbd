package calendar

import (
	"context"
)

type Service struct {
	repo        Repository
	permissions *PermissionService
}

func NewService(
	repo Repository,
	permissions *PermissionService,
) *Service {
	return &Service{
		repo:        repo,
		permissions: permissions,
	}
}

func (service *Service) Create(ctx context.Context, workerID string) (Calendar, error) {
	calendar := Calendar{ID: workerID}

	if err := service.repo.Create(ctx, calendar); err != nil {
		return Calendar{}, err
	}

	return calendar, service.permissions.NewCalendarPermissions(calendar)
}

func (service *Service) FindByID(ctx context.Context, ID string) (Calendar, error) {
	return service.repo.FindByID(ctx, ID)
}

func (service *Service) Delete(ctx context.Context, calendar Calendar) error {
	return service.repo.Delete(ctx, calendar)
}
