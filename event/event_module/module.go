package event_module

import (
	"app/event"
	"app/event/event_adapter"
	"app/event/event_http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var (
	ServiceModule = fx.Module(
		"event",

		fx.Provide(event_adapter.NewPostgresRepository, fx.Private),
		fx.Provide(func(repo *event_adapter.PostgresRepository) event.Repository { return repo }, fx.Private),

		fx.Provide(event.NewService),
	)

	HTTPModule = fx.Module(
		"event http",

		ServiceModule,

		fx.Provide(event_http.NewController, fx.Private),

		fx.Invoke(func(app *fiber.App, eventController *event_http.Controller) {
			eventController.Register(app)
		}),
	)
)
