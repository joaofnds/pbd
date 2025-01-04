package calendar_adapter

import (
	"app/calendar"
	"context"
	"errors"

	"gorm.io/gorm"
)

var _ calendar.Repository = (*PostgresRepository)(nil)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (repository *PostgresRepository) Create(ctx context.Context, cal calendar.Calendar) error {
	err := repository.db.
		WithContext(ctx).
		Exec("INSERT INTO calendars(id) VALUES(?)", cal.ID)

	return gormErr(err)
}

func (p *PostgresRepository) Delete(ctx context.Context, cal calendar.Calendar) error {
	err := p.db.WithContext(ctx).Exec("DELETE FROM calendars WHERE id = ?", cal.ID)

	return gormErr(err)
}

func (p *PostgresRepository) FindByID(ctx context.Context, id string) (calendar.Calendar, error) {
	var cal calendar.Calendar
	return cal, gormErr(p.db.WithContext(ctx).First(&cal, "id = ?", id))
}

func gormErr(result *gorm.DB) error {
	switch {
	case result.Error == nil:
		return nil
	case errors.Is(result.Error, gorm.ErrRecordNotFound):
		return calendar.ErrNotFound
	default:
		return calendar.ErrRepository
	}
}
