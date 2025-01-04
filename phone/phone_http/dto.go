package phone_http

import "app/phone"

type CreateBody struct {
	CountryCode string `json:"country_code" validate:"required,max=3,min=2"`
	AreaCode    string `json:"area_code" validate:"required,max=3,min=2"`
	Number      string `json:"number" validate:"required,max=9,min=8"`
}

func (dto *CreateBody) ToCreateDTO(userID string) phone.CreateDTO {
	return phone.CreateDTO{
		UserID:      userID,
		CountryCode: dto.CountryCode,
		AreaCode:    dto.AreaCode,
		Number:      dto.Number,
	}
}
