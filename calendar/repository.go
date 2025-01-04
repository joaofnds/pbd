package calendar

import "context"

type Repository interface {
	Create(context.Context, Calendar) error
	FindByID(context.Context, string) (Calendar, error)
	Delete(context.Context, Calendar) error
}
