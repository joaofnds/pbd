package order

import (
	"app/adapter/time"
	"app/address"
	"app/booking"
	"app/calendar"
	"app/event"
	"app/internal/clock"
	"app/internal/id"
	"app/worker"
	"context"
	"fmt"
)

type Service struct {
	id          id.Generator
	clock       clock.Clock
	repository  Repository
	permissions *PermissionService

	addresses *address.Service
	booking   *booking.Service
	calendars *calendar.Service
	events    *event.Service
	workers   *worker.Service
}

func NewService(
	id id.Generator,
	clock *time.Clock,
	repository Repository,
	permissions *PermissionService,

	addresses *address.Service,
	booking *booking.Service,
	calendars *calendar.Service,
	events *event.Service,
	workers *worker.Service,
) *Service {
	return &Service{
		id:          id,
		clock:       clock,
		repository:  repository,
		permissions: permissions,

		addresses: addresses,
		booking:   booking,
		calendars: calendars,
		events:    events,
		workers:   workers,
	}
}

func (service *Service) Create(ctx context.Context, request Request) (Order, error) {
	addr, findAddrErr := service.addresses.Get(ctx, request.AddressID)
	if findAddrErr != nil {
		return Order{}, fmt.Errorf("%w: %s", ErrOrder, findAddrErr)
	}

	evt, findEventErr := service.events.FindByID(ctx, request.EventID)
	if findEventErr != nil {
		return Order{}, fmt.Errorf("%w: %s", ErrOrder, findEventErr)
	}

	price, err := service.CalculatePrice(ctx, request)
	if err != nil {
		return Order{}, err
	}

	newEvent, bookingErr := service.booking.Book(ctx, evt, booking.Request{
		StartsAt: request.StartsAt,
		EndsAt:   request.EndsAt,
	})
	if bookingErr != nil {
		return Order{}, fmt.Errorf("%w: %s", ErrOrder, bookingErr)
	}

	order := Order{
		ID: service.id.NewID(),

		Price:  price,
		Status: StatusCreated,

		AddressID:  addr.ID,
		EventID:    newEvent.ID,
		WorkerID:   request.WorkerID,
		CustomerID: request.CustomerID,

		CreatedAt: service.clock.Now(),
		UpdatedAt: service.clock.Now(),
	}
	if createErr := service.repository.Create(ctx, order); createErr != nil {
		return Order{}, createErr
	}

	return order, service.permissions.GrantNewOrderPermissions(order)
}

func (service *Service) List(ctx context.Context) ([]Order, error) {
	return service.repository.List(ctx)
}

func (service *Service) GetByID(ctx context.Context, orderID string) (Order, error) {
	return service.repository.GetByID(ctx, orderID)
}

func (service *Service) FindByCustomerID(ctx context.Context, customerID string) ([]Order, error) {
	return service.repository.FindByCustomerID(ctx, customerID)
}

func (service *Service) FindByWorkerID(ctx context.Context, workerID string) ([]Order, error) {
	return service.repository.FindByWorkerID(ctx, workerID)
}

func (service *Service) Update(ctx context.Context, order Order, dto UpdateDTO) error {
	err := service.repository.Update(ctx, order.ID, dto)
	if err != nil {
		return err
	}

	if dto.Status == StatusCompleted {
		return service.permissions.GrantCompletedOrderPermissions(order)
	}

	return nil
}

func (service *Service) Delete(ctx context.Context, orderID string) error {
	return service.repository.Delete(ctx, orderID)
}

func (service *Service) CalculatePrice(ctx context.Context, request Request) (int, error) {
	worker, err := service.workers.FindByID(ctx, request.WorkerID)
	if err != nil {
		return 0, err
	}

	durationInHours := int(request.EndsAt.Sub(request.StartsAt).Hours())
	return durationInHours * worker.HourlyRate, nil
}
