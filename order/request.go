package order

import "time"

type Request struct {
	AddressID  string    `json:"address_id"`
	EventID    string    `json:"event_id"`
	WorkerID   string    `json:"worker_id"`
	CustomerID string    `json:"customer_id"`
	StartsAt   time.Time `json:"starts_at"`
	EndsAt     time.Time `json:"ends_at"`
}
