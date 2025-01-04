package phone

import "time"

type Phone struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	CountryCode string    `json:"country_code"`
	AreaCode    string    `json:"area_code"`
	Number      string    `json:"number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
