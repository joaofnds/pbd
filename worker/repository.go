package worker

import "context"

type Repository interface {
	Create(context.Context, Worker) error
	FindByID(context.Context, string) (Worker, error)
}
