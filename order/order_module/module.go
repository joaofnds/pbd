package order_module

import (
	"app/adapter/time"
	"app/adapter/uuid"
	"app/internal/clock"
	"app/internal/id"
	"app/order"
	"app/order/order_adapter"
	"app/order/order_http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var (
	ServiceModule = fx.Module(
		"order",

		fx.Provide(order_adapter.NewPostgresRepository, fx.Private),
		fx.Provide(func(repo *order_adapter.PostgresRepository) order.Repository { return repo }, fx.Private),
		fx.Provide(func(gen *uuid.Generator) id.Generator { return gen }, fx.Private),
		fx.Provide(func(clock *time.Clock) clock.Clock { return clock }, fx.Private),

		fx.Provide(order.NewService),
		fx.Provide(order.NewPermissionService),
	)

	HTTPModule = fx.Module(
		"order http",

		ServiceModule,

		fx.Provide(order_http.NewController, fx.Private),

		fx.Invoke(func(app *fiber.App, orderController *order_http.Controller) {
			orderController.Register(app)
		}),
	)
)
