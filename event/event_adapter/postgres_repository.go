package event_adapter

import (
	"app/event"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

var _ event.Repository = (*PostgresRepository)(nil)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db.Table("events")}
}

func (repository *PostgresRepository) Create(ctx context.Context, evt event.Event) error {
	err := repository.db.
		WithContext(ctx).
		Create(&evt)

	return gormErr(err)
}

func (repository *PostgresRepository) List(ctx context.Context, dto event.ListDTO) ([]event.Event, error) {
	var events []event.Event
	err := repository.db.
		WithContext(ctx).
		Where("calendar_id = ?", dto.CalendarID).
		Where("starts_at >= ?", dto.StartsAt).
		Where("ends_at <= ?", dto.EndsAt).
		Find(&events)

	return events, gormErr(err)
}

func (repository *PostgresRepository) FindByID(ctx context.Context, eventID string) (event.Event, error) {
	var evt event.Event
	return evt, gormErr(repository.db.WithContext(ctx).First(&evt, "id = ?", eventID))
}

func (repository *PostgresRepository) FindIntersecting(ctx context.Context, dto event.ListDTO) ([]event.Event, error) {
	var events []event.Event
	err := repository.db.
		WithContext(ctx).
		Where("calendar_id = ?", dto.CalendarID).
		Where("starts_at < ?", dto.EndsAt).
		Where("ends_at > ?", dto.StartsAt).
		Find(&events)

	return events, gormErr(err)
}

func (repository *PostgresRepository) Update(ctx context.Context, evt event.Event, dto event.UpdateDTO) error {
	err := repository.db.
		WithContext(ctx).
		Where("id = ?", evt.ID).
		Updates(dto)

	return gormErr(err)
}

func (repository *PostgresRepository) Delete(ctx context.Context, evt event.Event) error {
	err := repository.db.WithContext(ctx).Delete(&evt)

	return gormErr(err)
}

func (repository *PostgresRepository) IsTimeSlotTaken(
	ctx context.Context,
	calendarID string,
	startsAt time.Time,
	endsAt time.Time,
) bool {
	result := repository.db.
		WithContext(ctx).
		First(
			&event.Event{},
			"status in (?) AND calendar_id = ? AND starts_at < ? AND ends_at > ?",
			event.NonOverlappableStatuses,
			calendarID,
			endsAt,
			startsAt,
		)
	return !errors.Is(result.Error, gorm.ErrRecordNotFound)
}

func gormErr(result *gorm.DB) error {
	switch {
	case result.Error == nil:
		return nil
	case errors.Is(result.Error, gorm.ErrRecordNotFound):
		return event.ErrNotFound
	default:
		return event.ErrRepository
	}
}
