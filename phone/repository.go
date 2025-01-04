package phone

import "context"

type Repository interface {
	Create(context.Context, Phone) error
	Get(context.Context, string) (Phone, error)
	Delete(context.Context, string) error

	FindByUserID(context.Context, string) ([]Phone, error)
}
