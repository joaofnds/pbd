package customer

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

func (service *Service) Create(ctx context.Context, userID string) (Customer, error) {
	customer := Customer{
		ID:        userID,
		CreatedAt: service.clock.Now(),
		UpdatedAt: service.clock.Now(),
	}

	if err := service.repo.Create(ctx, customer); err != nil {
		return Customer{}, err
	}

	return customer, service.permissions.GrantNewCustomerPermissions(customer)
}

func (service *Service) FindByID(ctx context.Context, id string) (Customer, error) {
	return service.repo.FindByID(ctx, id)
}
