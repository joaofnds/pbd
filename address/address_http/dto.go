package address_http

import "app/address"

type CreateBody struct {
	Street       string `json:"street" validate:"required,max=255,min=1"`
	Number       string `json:"number" validate:"max=50"`
	Complement   string `json:"complement" validate:"max=255"`
	Neighborhood string `json:"neighborhood" validate:"required,max=255,min=1"`
	City         string `json:"city" validate:"required,max=255,min=1"`
	State        string `json:"state" validate:"required,len=2"`
	ZipCode      string `json:"zip_code" validate:"required,min=5,max=9"`
	Country      string `json:"country" validate:"required,max=255,min=1"`
}

func (dto *CreateBody) ToCreateDTO(customerID string) address.CreateDTO {
	return address.CreateDTO{
		CustomerID:   customerID,
		Street:       dto.Street,
		Number:       dto.Number,
		Complement:   dto.Complement,
		Neighborhood: dto.Neighborhood,
		City:         dto.City,
		State:        dto.State,
		ZipCode:      dto.ZipCode,
		Country:      dto.Country,
	}
}
