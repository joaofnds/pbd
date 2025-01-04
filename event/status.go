package event

const (
	StatusAvailable = "available"
	StatusBooked    = "booked"
	StatusCanceled  = "canceled"
)

var (
	OverlappableStatuses    = []string{StatusCanceled}
	NonOverlappableStatuses = []string{StatusAvailable, StatusBooked}
)
