package customer

import (
	"app/authn"
	"context"
)

type RegistrationService struct {
	auth      *authn.Service
	customers *Service
}

func NewRegistrationService(
	auth *authn.Service,
	customers *Service,
) *RegistrationService {
	return &RegistrationService{
		auth:      auth,
		customers: customers,
	}
}

func (service *RegistrationService) Register(
	ctx context.Context,
	email string,
	password string,
) (Customer, error) {
	usr, registerErr := service.auth.RegisterUser(ctx, email, password)
	if registerErr != nil {
		return Customer{}, registerErr
	}

	return service.customers.Create(ctx, usr.ID)
}
