package phone

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

func (service *Service) Create(ctx context.Context, dto CreateDTO) (Phone, error) {
	phone := Phone{
		ID:        service.id.NewID(),
		CreatedAt: service.clock.Now(),
		UpdatedAt: service.clock.Now(),

		UserID:      dto.UserID,
		CountryCode: dto.CountryCode,
		AreaCode:    dto.AreaCode,
		Number:      dto.Number,
	}

	if err := service.repository.Create(ctx, phone); err != nil {
		return Phone{}, err
	}

	return phone, service.permissions.GrantNewPhonePermissions(phone)
}

func (service *Service) Get(ctx context.Context, id string) (Phone, error) {
	return service.repository.Get(ctx, id)
}

func (service *Service) Delete(ctx context.Context, id string) error {
	return service.repository.Delete(ctx, id)
}

func (service *Service) FindByUserID(ctx context.Context, userID string) ([]Phone, error) {
	return service.repository.FindByUserID(ctx, userID)
}
