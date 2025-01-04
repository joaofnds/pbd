package order_http

import "time"

type CreatePayload struct {
	AddressID  string    `json:"address_id" validate:"required,uuid4"`
	EventID    string    `json:"event_id" validate:"required,uuid4"`
	WorkerID   string    `json:"worker_id" validate:"required,uuid4"`
	CustomerID string    `json:"customer_id" validate:"required,uuid4"`
	StartsAt   time.Time `json:"starts_at" validate:"required"`
	EndsAt     time.Time `json:"ends_at" validate:"required"`
}

type UpdatePayload struct {
	Status string `json:"status" validate:"required,oneof=completed"`
}
