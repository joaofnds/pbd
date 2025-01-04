package worker

import (
	"time"
)

type Worker struct {
	ID         string    `json:"id"`
	HourlyRate int       `json:"hourly_rate"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
