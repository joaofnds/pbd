package booking_module

import (
	"app/booking"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"booking",

	fx.Provide(booking.NewService),
)
