package worker

import (
	"app/authn"
	"app/calendar"
	"context"
)

type RegistrationService struct {
	auth      *authn.Service
	workers   *Service
	calendars *calendar.Service
}

func NewRegistrationService(
	auth *authn.Service,
	workers *Service,
	calendars *calendar.Service,
) *RegistrationService {
	return &RegistrationService{
		auth:      auth,
		workers:   workers,
		calendars: calendars,
	}
}

func (service *RegistrationService) Register(ctx context.Context, dto RegisterDTO) (Worker, error) {
	usr, registerErr := service.auth.RegisterUser(ctx, dto.Email, dto.Password)
	if registerErr != nil {
		return Worker{}, registerErr
	}

	worker, createWorkerErr := service.workers.Create(ctx, usr.ID, dto.HourlyRate)
	if createWorkerErr != nil {
		return Worker{}, createWorkerErr
	}

	_, createCalendarErr := service.calendars.Create(ctx, worker.ID)
	if createCalendarErr != nil {
		return Worker{}, createCalendarErr
	}

	return worker, nil
}
