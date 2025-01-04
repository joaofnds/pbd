package address_adapter

import (
	"app/address"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var _ address.Repository = (*PostgresRepository)(nil)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db.Table("addresses")}
}

func (repository *PostgresRepository) Create(ctx context.Context, addressToCreate address.Address) error {
	return gormErr(repository.db.WithContext(ctx).Create(&addressToCreate))
}

func (repository *PostgresRepository) Get(ctx context.Context, id string) (address.Address, error) {
	var address address.Address
	result := repository.db.
		WithContext(ctx).
		First(&address, "id = ?", id)
	return address, gormErr(result)
}

func (repository *PostgresRepository) Delete(ctx context.Context, id string) error {
	result := repository.db.
		WithContext(ctx).
		Delete(&address.Address{}, "id = ?", id)
	return gormErr(result)
}

func (repository *PostgresRepository) FindByCustomerID(ctx context.Context, customerID string) ([]address.Address, error) {
	var addresses []address.Address
	result := repository.db.
		WithContext(ctx).
		Where("customer_id = ?", customerID).
		Find(&addresses)
	return addresses, gormErr(result)
}

func gormErr(result *gorm.DB) error {
	switch {
	case result.Error == nil:
		return nil
	case errors.Is(result.Error, gorm.ErrRecordNotFound):
		return address.ErrNotFound
	default:
		return fmt.Errorf("%w: %v", address.ErrRepository, result.Error)
	}
}
