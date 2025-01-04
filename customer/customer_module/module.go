package customer_module

import (
	"app/customer"
	"app/customer/customer_adapter"
	"app/customer/customer_http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var (
	ServiceModule = fx.Module(
		"customer",

		fx.Provide(customer_adapter.NewPostgresRepository, fx.Private),
		fx.Provide(func(repo *customer_adapter.PostgresRepository) customer.Repository { return repo }, fx.Private),

		fx.Provide(customer.NewService),
		fx.Provide(customer.NewRegistrationService),
		fx.Provide(customer.NewPermissionService),
	)

	HTTPModule = fx.Module(
		"customer http",

		ServiceModule,

		fx.Provide(customer_http.NewController, fx.Private),

		fx.Invoke(func(app *fiber.App, customerController *customer_http.Controller) {
			customerController.Register(app)
		}),
	)
)
