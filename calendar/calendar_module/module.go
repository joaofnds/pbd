package calendar_module

import (
	"app/adapter/time"
	"app/adapter/uuid"
	"app/calendar"
	"app/calendar/calendar_adapter"
	"app/calendar/calendar_http"
	"app/internal/clock"
	"app/internal/id"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var (
	ServiceModule = fx.Module(
		"calendar",

		fx.Provide(calendar_adapter.NewPostgresRepository, fx.Private),
		fx.Provide(func(repo *calendar_adapter.PostgresRepository) calendar.Repository { return repo }, fx.Private),
		fx.Provide(func(gen *uuid.Generator) id.Generator { return gen }, fx.Private),
		fx.Provide(func(clock *time.Clock) clock.Clock { return clock }, fx.Private),

		fx.Provide(calendar.NewService),
		fx.Provide(calendar.NewPermissionService),
	)

	HTTPModule = fx.Module(
		"calendar http",
		ServiceModule,

		fx.Provide(calendar_http.NewController, fx.Private),

		fx.Invoke(func(app *fiber.App, calendarController *calendar_http.Controller) {
			calendarController.Register(app)
		}),
	)
)
