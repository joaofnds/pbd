package order

import (
	"time"
)

type Order struct {
	ID string `json:"id"`

	Price  int    `json:"price"`
	Status string `json:"status"`

	AddressID  string `json:"address_id"`
	EventID    string `json:"event_id"`
	WorkerID   string `json:"worker_id"`
	CustomerID string `json:"customer_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
