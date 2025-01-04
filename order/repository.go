package order

import "context"

type Repository interface {
	Create(context.Context, Order) error
	List(context.Context) ([]Order, error)
	GetByID(context.Context, string) (Order, error)
	FindByCustomerID(context.Context, string) ([]Order, error)
	FindByWorkerID(context.Context, string) ([]Order, error)
	Update(context.Context, string, UpdateDTO) error
	Delete(context.Context, string) error
}
