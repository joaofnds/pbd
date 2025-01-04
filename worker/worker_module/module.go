package worker_module

import (
	"app/worker"
	"app/worker/worker_adapter"
	"app/worker/worker_http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var (
	ServiceModule = fx.Module(
		"worker",

		fx.Provide(worker_adapter.NewPostgresRepository, fx.Private),
		fx.Provide(func(repo *worker_adapter.PostgresRepository) worker.Repository { return repo }, fx.Private),

		fx.Provide(worker.NewService),
		fx.Provide(worker.NewRegistrationService),
		fx.Provide(worker.NewPermissionServiceService),
	)

	HTTPModule = fx.Module(
		"worker http",

		ServiceModule,

		fx.Provide(worker_http.NewController, fx.Private),

		fx.Invoke(func(app *fiber.App, workerController *worker_http.Controller) {
			workerController.Register(app)
		}),
	)
)
