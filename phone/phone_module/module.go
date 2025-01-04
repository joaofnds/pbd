package phone_module

import (
	"app/adapter/time"
	"app/adapter/uuid"
	"app/internal/clock"
	"app/internal/id"
	"app/phone"
	"app/phone/phone_adapter"
	"app/phone/phone_http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var (
	ServiceModule = fx.Module(
		"phone",

		fx.Provide(phone_adapter.NewPostgresRepository, fx.Private),
		fx.Provide(func(repo *phone_adapter.PostgresRepository) phone.Repository { return repo }, fx.Private),
		fx.Provide(func(gen *uuid.Generator) id.Generator { return gen }, fx.Private),
		fx.Provide(func(clock *time.Clock) clock.Clock { return clock }, fx.Private),

		fx.Provide(phone.NewService),
		fx.Provide(phone.NewPermissionService),
	)

	HTTPModule = fx.Module(
		"phone http",

		ServiceModule,

		fx.Provide(phone_http.NewController, fx.Private),

		fx.Invoke(func(app *fiber.App, phoneController *phone_http.Controller) {
			phoneController.Register(app)
		}),
	)
)
