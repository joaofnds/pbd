package worker_adapter

import (
	"app/worker"
	"context"
	"errors"

	"gorm.io/gorm"
)

var _ worker.Repository = (*PostgresRepository)(nil)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db}
}

func (repository *PostgresRepository) Create(ctx context.Context, newWorker worker.Worker) error {
	err := repository.db.
		WithContext(ctx).
		Exec("INSERT INTO workers(id, hourly_rate, created_at, updated_at) VALUES(?, ?, ?, ?)",
			newWorker.ID,
			newWorker.HourlyRate,
			newWorker.CreatedAt,
			newWorker.UpdatedAt,
		)

	return gormErr(err)
}

func (repository *PostgresRepository) FindByID(ctx context.Context, id string) (worker.Worker, error) {
	var foundWorker worker.Worker
	return foundWorker, gormErr(repository.db.WithContext(ctx).First(&foundWorker, "id = ?", id))
}

func gormErr(result *gorm.DB) error {
	switch {
	case result.Error == nil:
		return nil
	case errors.Is(result.Error, gorm.ErrRecordNotFound):
		return worker.ErrNotFound
	default:
		return worker.ErrRepository
	}
}
