package address

import "time"

type Address struct {
	ID           string    `json:"id"`
	CustomerID   string    `json:"customer_id"`
	Street       string    `json:"street"`
	Number       string    `json:"number,omitempty"`
	Complement   string    `json:"complement,omitempty"`
	Neighborhood string    `json:"neighborhood"`
	City         string    `json:"city"`
	State        string    `json:"state"`
	ZipCode      string    `json:"zip_code"`
	Country      string    `json:"country"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
