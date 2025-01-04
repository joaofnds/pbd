package address

import (
	"app/internal/clock"
	"app/internal/id"
	"context"
)

type Service struct {
	id          id.Generator
	clock       clock.Clock
	repository  Repository
	permissions *PermissionService
}

func NewService(
	id id.Generator,
	clock clock.Clock,
	repository Repository,
	permissions *PermissionService,
) *Service {
	return &Service{
		id:          id,
		clock:       clock,
		repository:  repository,
		permissions: permissions,
	}
}

func (service *Service) Create(ctx context.Context, dto CreateDTO) (Address, error) {
	address := Address{
		ID:        service.id.NewID(),
		CreatedAt: service.clock.Now(),
		UpdatedAt: service.clock.Now(),

		CustomerID:   dto.CustomerID,
		Street:       dto.Street,
		Number:       dto.Number,
		Complement:   dto.Complement,
		Neighborhood: dto.Neighborhood,
		City:         dto.City,
		State:        dto.State,
		ZipCode:      dto.ZipCode,
		Country:      dto.Country,
	}

	if err := service.repository.Create(ctx, address); err != nil {
		return Address{}, err
	}

	return address, service.permissions.GrantNewAddressPermissions(address)
}

func (service *Service) Get(ctx context.Context, id string) (Address, error) {
	return service.repository.Get(ctx, id)
}

func (service *Service) Delete(ctx context.Context, id string) error {
	return service.repository.Delete(ctx, id)
}

func (service *Service) FindByCustomerID(ctx context.Context, customerID string) ([]Address, error) {
	return service.repository.FindByCustomerID(ctx, customerID)
}
