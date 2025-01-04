package customer

import "context"

type Repository interface {
	Create(context.Context, Customer) error
	FindByID(context.Context, string) (Customer, error)
}
