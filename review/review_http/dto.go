package review_http

type CreateReviewDTO struct {
	Rating  int    `json:"user_id" validate:"required,gte=1,lte=100"`
	Comment string `json:"comment"`
}
