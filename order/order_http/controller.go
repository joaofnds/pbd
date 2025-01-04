package order_http

import (
	"app/adapter/validation"
	"app/authn/authn_http"
	"app/authz/authz_http"
	"app/order"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Controller struct {
	logger    *zap.Logger
	validator *validator.Validate
	authn     *authn_http.AuthMiddleware
	authz     *authz_http.Middleware
	service   *order.Service
}

func NewController(
	logger *zap.Logger,
	validator *validator.Validate,
	authn *authn_http.AuthMiddleware,
	authz *authz_http.Middleware,
	service *order.Service,
) *Controller {
	return &Controller{
		logger:    logger.Named("order-controller"),
		validator: validator,
		authn:     authn,
		authz:     authz,
		service:   service,
	}
}

func (controller *Controller) Register(app *fiber.App) {
	orders := app.Group("/orders", controller.authn.RequireUser)
	orders.
		Post("/", controller.Create).
		Get("/", controller.List)

	orders.Group("/:orderID", controller.middlewareGetOrder).
		Get(
			"/",
			controller.authz.RequireParamPermission("order:orderID", order.PermOrderRead),
			controller.Get,
		).
		Patch(
			"/",
			controller.authz.RequireParamPermission("order:orderID", order.PermOrderUpdate),
			controller.Update,
		).
		Delete(
			"/",
			controller.authz.RequireParamPermission("order:orderID", order.PermOrderDelete),
			controller.Delete,
		)
}

func (controller *Controller) Create(ctx *fiber.Ctx) error {
	var body CreatePayload
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	if err := controller.validator.Struct(body); err != nil {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(fiber.Map{"errors": validation.ErrorMessages(err)})
	}

	ord, createErr := controller.service.Create(ctx.UserContext(), order.Request{
		AddressID:  body.AddressID,
		EventID:    body.EventID,
		WorkerID:   body.WorkerID,
		CustomerID: body.CustomerID,
		StartsAt:   body.StartsAt.UTC(),
		EndsAt:     body.EndsAt.UTC(),
	})

	switch {
	case createErr == nil:
		return ctx.Status(http.StatusCreated).JSON(ord)
	case errors.Is(createErr, order.ErrOrder):
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": createErr.Error()})
	default:
		controller.logger.Error("failed to create order", zap.Error(createErr))
		return ctx.SendStatus(http.StatusInternalServerError)
	}
}

func (controller *Controller) List(ctx *fiber.Ctx) error {
	orders, err := controller.service.List(ctx.UserContext())
	if err != nil {
		controller.logger.Error("failed to list orders", zap.Error(err))
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.JSON(orders)
}

func (controller *Controller) Get(ctx *fiber.Ctx) error {
	ord, err := controller.service.GetByID(ctx.UserContext(), ctx.Params("orderID"))
	if err != nil {
		controller.logger.Error("failed to get order", zap.Error(err))
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.JSON(ord)
}

func (controller *Controller) Update(ctx *fiber.Ctx) error {
	var body UpdatePayload
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	if err := controller.validator.Struct(body); err != nil {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(fiber.Map{"errors": validation.ErrorMessages(err)})
	}

	ord := ctx.Locals("order").(order.Order)

	if err := controller.service.Update(ctx.UserContext(), ord, order.UpdateDTO{
		Status: body.Status,
	}); err != nil {
		controller.logger.Error("failed to update order", zap.Error(err))
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.SendStatus(http.StatusOK)
}

func (controller *Controller) Delete(ctx *fiber.Ctx) error {
	err := controller.service.Delete(ctx.UserContext(), ctx.Params("orderID"))
	if err != nil {
		controller.logger.Error("failed to delete order", zap.Error(err))
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.SendStatus(http.StatusNoContent)
}

func (controller *Controller) middlewareGetOrder(ctx *fiber.Ctx) error {
	orderID := ctx.Params("orderID")
	if orderID == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	ord, err := controller.service.GetByID(ctx.UserContext(), orderID)
	switch {
	case err == nil:
		ctx.Locals("order", ord)
		return ctx.Next()
	case errors.Is(err, order.ErrNotFound):
		return ctx.SendStatus(http.StatusNotFound)
	default:
		return ctx.SendStatus(http.StatusInternalServerError)
	}
}
