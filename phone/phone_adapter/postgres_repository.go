package phone_adapter

import (
	"app/phone"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var _ phone.Repository = (*PostgresRepository)(nil)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db.Table("phones")}
}

func (repository *PostgresRepository) Create(ctx context.Context, phoneToCreate phone.Phone) error {
	return gormErr(repository.db.WithContext(ctx).Create(&phoneToCreate))
}

func (repository *PostgresRepository) Get(ctx context.Context, id string) (phone.Phone, error) {
	var phone phone.Phone
	result := repository.db.
		WithContext(ctx).
		First(&phone, "id = ?", id)
	return phone, gormErr(result)
}

func (repository *PostgresRepository) Delete(ctx context.Context, id string) error {
	result := repository.db.
		WithContext(ctx).
		Delete(&phone.Phone{}, "id = ?", id)
	return gormErr(result)
}

func (repository *PostgresRepository) FindByUserID(ctx context.Context, userID string) ([]phone.Phone, error) {
	var phones []phone.Phone
	result := repository.db.
		WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&phones)
	return phones, gormErr(result)
}

func gormErr(result *gorm.DB) error {
	switch {
	case result.Error == nil:
		return nil
	case errors.Is(result.Error, gorm.ErrRecordNotFound):
		return phone.ErrNotFound
	default:
		return fmt.Errorf("%w: %v", phone.ErrRepository, result.Error)
	}
}
