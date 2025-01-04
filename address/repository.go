package address

import "context"

type Repository interface {
	Create(context.Context, Address) error
	Get(context.Context, string) (Address, error)
	Delete(context.Context, string) error

	FindByCustomerID(context.Context, string) ([]Address, error)
}
