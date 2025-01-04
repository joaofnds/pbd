package review

import "context"

type Repository interface {
	Create(context.Context, Review) error
	GetByOrderID(context.Context, string) (Review, error)
	FindByWorkerID(context.Context, string) ([]Review, error)
	FindByCustomerID(context.Context, string) ([]Review, error)
	Delete(context.Context, Review) error
}
