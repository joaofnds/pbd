package phone

type CreateDTO struct {
	UserID      string `json:"user_id"`
	CountryCode string `json:"country_code"`
	AreaCode    string `json:"area_code"`
	Number      string `json:"number"`
}
