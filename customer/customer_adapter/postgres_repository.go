package customer_adapter

import (
	"app/customer"
	"context"
	"errors"

	"gorm.io/gorm"
)

var _ customer.Repository = (*PostgresRepository)(nil)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db}
}

func (repository *PostgresRepository) Create(ctx context.Context, newCustomer customer.Customer) error {
	err := repository.db.
		WithContext(ctx).
		Exec("INSERT INTO customers(id, created_at, updated_at) VALUES(?, ?, ?)",
			newCustomer.ID,
			newCustomer.CreatedAt,
			newCustomer.UpdatedAt,
		)

	return gormErr(err)
}

func (repository *PostgresRepository) FindByID(ctx context.Context, id string) (customer.Customer, error) {
	var foundCustomer customer.Customer
	return foundCustomer, gormErr(repository.db.WithContext(ctx).First(&foundCustomer, "id = ?", id))
}

func gormErr(result *gorm.DB) error {
	switch {
	case result.Error == nil:
		return nil
	case errors.Is(result.Error, gorm.ErrRecordNotFound):
		return customer.ErrNotFound
	default:
		return customer.ErrRepository
	}
}
