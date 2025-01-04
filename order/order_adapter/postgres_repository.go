package order_adapter

import (
	"app/order"
	"context"
	"errors"

	"gorm.io/gorm"
)

var _ order.Repository = (*PostgresRepository)(nil)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db.Table("orders")}
}

func (repository *PostgresRepository) Create(ctx context.Context, ord order.Order) error {
	err := repository.db.
		WithContext(ctx).
		Create(&ord)

	return gormErr(err)
}

func (repository *PostgresRepository) GetByID(ctx context.Context, orderID string) (order.Order, error) {
	var ord order.Order
	return ord, gormErr(repository.db.WithContext(ctx).First(&ord, "id = ?", orderID))
}

func (repository *PostgresRepository) FindByCustomerID(ctx context.Context, customerID string) ([]order.Order, error) {
	var orders []order.Order
	return orders, gormErr(repository.db.WithContext(ctx).Find(&orders, "customer_id = ?", customerID))
}

func (repository *PostgresRepository) FindByWorkerID(ctx context.Context, workerID string) ([]order.Order, error) {
	var orders []order.Order
	return orders, gormErr(repository.db.WithContext(ctx).Find(&orders, "worker_id = ?", workerID))
}

func (repository *PostgresRepository) List(ctx context.Context) ([]order.Order, error) {
	var orders []order.Order
	return orders, gormErr(repository.db.WithContext(ctx).Find(&orders))
}

func (repository *PostgresRepository) Update(ctx context.Context, orderID string, dto order.UpdateDTO) error {
	err := repository.db.WithContext(ctx).Model(&order.Order{}).Where("id = ?", orderID).Updates(dto)
	return gormErr(err)
}

func (repository *PostgresRepository) Delete(ctx context.Context, orderID string) error {
	err := repository.db.WithContext(ctx).Delete(&order.Order{}, "id = ?", orderID)
	return gormErr(err)
}

func gormErr(result *gorm.DB) error {
	switch {
	case result.Error == nil:
		return nil
	case errors.Is(result.Error, gorm.ErrRecordNotFound):
		return order.ErrNotFound
	default:
		return order.ErrRepository
	}
}
