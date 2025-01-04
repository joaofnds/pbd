package address_module

import (
	"app/adapter/time"
	"app/adapter/uuid"
	"app/address"
	"app/address/address_adapter"
	"app/address/address_http"
	"app/internal/clock"
	"app/internal/id"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var (
	ServiceModule = fx.Module(
		"address",

		fx.Provide(address_adapter.NewPostgresRepository, fx.Private),
		fx.Provide(func(repo *address_adapter.PostgresRepository) address.Repository { return repo }, fx.Private),
		fx.Provide(func(gen *uuid.Generator) id.Generator { return gen }, fx.Private),
		fx.Provide(func(clock *time.Clock) clock.Clock { return clock }, fx.Private),

		fx.Provide(address.NewService),
		fx.Provide(address.NewPermissionService),
	)

	HTTPModule = fx.Module(
		"address http",

		ServiceModule,

		fx.Provide(address_http.NewController, fx.Private),

		fx.Invoke(func(app *fiber.App, addressController *address_http.Controller) {
			addressController.Register(app)
		}),
	)
)
