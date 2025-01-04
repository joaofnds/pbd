package review_module

import (
	"app/adapter/time"
	"app/adapter/uuid"
	"app/internal/clock"
	"app/internal/id"
	"app/review"
	"app/review/review_adapter"
	"app/review/review_http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var (
	ServiceModule = fx.Module(
		"review",

		fx.Provide(review_adapter.NewPostgresRepository, fx.Private),
		fx.Provide(func(repo *review_adapter.PostgresRepository) review.Repository { return repo }, fx.Private),
		fx.Provide(func(gen *uuid.Generator) id.Generator { return gen }, fx.Private),
		fx.Provide(func(clock *time.Clock) clock.Clock { return clock }, fx.Private),

		fx.Provide(review.NewService),
	)

	HTTPModule = fx.Module(
		"review http",

		ServiceModule,

		fx.Provide(review_http.NewController, fx.Private),

		fx.Invoke(func(app *fiber.App, reviewController *review_http.Controller) {
			reviewController.Register(app)
		}),
	)
)
