package worker_http

import (
	"errors"
	"net/http"

	"app/adapter/validation"
	"app/authn/authn_http"
	"app/authz/authz_http"
	"app/order"
	"app/review"
	"app/worker"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func NewController(
	validator *validator.Validate,
	authn *authn_http.AuthMiddleware,
	authz *authz_http.Middleware,
	workers *worker.Service,
	orders *order.Service,
	reviews *review.Service,
	registration *worker.RegistrationService,
) *Controller {
	return &Controller{
		validator:    validator,
		authn:        authn,
		authz:        authz,
		workers:      workers,
		orders:       orders,
		reviews:      reviews,
		registration: registration,
	}
}

type Controller struct {
	validator    *validator.Validate
	authn        *authn_http.AuthMiddleware
	authz        *authz_http.Middleware
	workers      *worker.Service
	orders       *order.Service
	reviews      *review.Service
	registration *worker.RegistrationService
}

func (controller *Controller) Register(app *fiber.App) {
	workers := app.Group("/workers")
	workers.Post("/", controller.Create)

	worker := app.Group(
		"/workers/:workerID",
		controller.authn.RequireUser,
		controller.middlewareGetWorker,
	)
	worker.Get(
		"/",
		controller.authz.RequireParamPermission("worker:workerID", "read"),
		controller.Get,
	)

	orders := worker.Group("/orders")
	orders.Get(
		"/",
		controller.authz.RequireParamPermission("worker:workerID", "order:list"),
		controller.ListOrders,
	)

	reviews := worker.Group("/reviews")
	reviews.Get("/", controller.ListReviews)
}

func (controller *Controller) Create(ctx *fiber.Ctx) error {
	var body CreateWorkerDTO
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	if err := controller.validator.Struct(body); err != nil {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(fiber.Map{"errors": validation.ErrorMessages(err)})
	}

	createdWorker, err := controller.registration.Register(ctx.UserContext(), body.ToRegisterDTO())
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.Status(http.StatusCreated).JSON(createdWorker)
}

func (controller *Controller) Get(ctx *fiber.Ctx) error {
	return ctx.JSON(ctx.Locals("worker"))
}

func (controller *Controller) ListOrders(ctx *fiber.Ctx) error {
	worker := ctx.Locals("worker").(worker.Worker)
	orders, err := controller.orders.FindByWorkerID(ctx.UserContext(), worker.ID)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.JSON(orders)
}

func (controller *Controller) ListReviews(ctx *fiber.Ctx) error {
	worker := ctx.Locals("worker").(worker.Worker)

	reviews, err := controller.reviews.FindByWorkerID(ctx.UserContext(), worker.ID)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.JSON(reviews)
}

func (controller *Controller) middlewareGetWorker(ctx *fiber.Ctx) error {
	name := ctx.Params("workerID")
	if name == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	workerFound, err := controller.workers.FindByID(ctx.UserContext(), name)
	if err != nil {
		if errors.Is(err, worker.ErrNotFound) {
			return ctx.SendStatus(http.StatusNotFound)
		} else {
			return ctx.SendStatus(http.StatusInternalServerError)
		}
	}

	ctx.Locals("worker", workerFound)
	return ctx.Next()
}
