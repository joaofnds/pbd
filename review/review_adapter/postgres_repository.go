package review_adapter

import (
	"app/review"
	"context"
	"errors"

	"gorm.io/gorm"
)

var _ review.Repository = (*PostgresRepository)(nil)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (repository *PostgresRepository) Create(ctx context.Context, reviewToCreate review.Review) error {
	result := repository.db.WithContext(ctx).Create(&reviewToCreate)
	return gormErr(result)
}

func (repository *PostgresRepository) GetByOrderID(ctx context.Context, orderID string) (review.Review, error) {
	var review review.Review
	result := repository.db.
		WithContext(ctx).
		Where("order_id = ?", orderID).
		First(&review)

	return review, gormErr(result)
}

func (repository *PostgresRepository) FindByCustomerID(ctx context.Context, customerID string) ([]review.Review, error) {
	query := `
    SELECT reviews.*
    FROM reviews
    INNER JOIN orders ON reviews.order_id = orders.id
    WHERE orders.customer_id = ?
    ORDER BY reviews.created_at DESC
  `

	var reviews []review.Review
	result := repository.db.WithContext(ctx).Raw(query, customerID).Scan(&reviews)
	return reviews, gormErr(result)
}

func (repository *PostgresRepository) FindByWorkerID(ctx context.Context, workerID string) ([]review.Review, error) {
	query := `
    SELECT reviews.*
    FROM reviews
    INNER JOIN orders ON reviews.order_id = orders.id
    WHERE orders.worker_id = ?
    ORDER BY reviews.created_at DESC
  `

	var reviews []review.Review
	result := repository.db.WithContext(ctx).Raw(query, workerID).Scan(&reviews)
	return reviews, gormErr(result)
}

func (p *PostgresRepository) Delete(ctx context.Context, reviewToDelete review.Review) error {
	result := p.db.WithContext(ctx).Exec("DELETE FROM reviews WHERE id = ?", reviewToDelete.ID)
	return gormErr(result)
}

func gormErr(result *gorm.DB) error {
	switch {
	case result.Error == nil:
		return nil
	case errors.Is(result.Error, gorm.ErrRecordNotFound):
		return review.ErrNotFound
	default:
		return review.ErrRepository
	}
}
