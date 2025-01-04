package worker_http

import "app/worker"

type CreateWorkerDTO struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required"`
	HourlyRate int    `json:"hourly_rate" validate:"required,gte=100"`
}

func (dto *CreateWorkerDTO) ToRegisterDTO() worker.RegisterDTO {
	return worker.RegisterDTO{
		Email:      dto.Email,
		Password:   dto.Password,
		HourlyRate: dto.HourlyRate,
	}
}
