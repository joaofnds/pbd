package review_http

import (
	"errors"
	"net/http"

	"app/adapter/validation"
	"app/authn/authn_http"
	"app/authz/authz_http"
	"app/order"
	"app/review"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func NewController(
	validator *validator.Validate,
	authn *authn_http.AuthMiddleware,
	authz *authz_http.Middleware,
	orders *order.Service,
	reviews *review.Service,
) *Controller {
	return &Controller{
		validator: validator,
		authn:     authn,
		authz:     authz,
		orders:    orders,
		reviews:   reviews,
	}
}

type Controller struct {
	validator *validator.Validate
	authn     *authn_http.AuthMiddleware
	authz     *authz_http.Middleware
	orders    *order.Service
	reviews   *review.Service
}

func (controller *Controller) Register(app *fiber.App) {
	reviews := app.Group(
		"/orders/:orderID/reviews",
		controller.middlewareGetOrder,
		controller.authn.RequireUser,
	)
	reviews.Get(
		"/",
		controller.authz.RequireParamPermission("order:orderID", order.PermReviewRead),
		controller.Get,
	)
	reviews.Post(
		"/",
		controller.authz.RequireParamPermission("order:orderID", order.PermReviewCreate),
		controller.Create,
	)
	reviews.Delete(
		"/",
		controller.authz.RequireParamPermission("order:orderID", order.PermReviewDelete),
		controller.Delete,
	)
}

func (controller *Controller) Create(ctx *fiber.Ctx) error {
	order := ctx.Locals("order").(order.Order)
	var body CreateReviewDTO
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	if err := controller.validator.Struct(body); err != nil {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(fiber.Map{"errors": validation.ErrorMessages(err)})
	}

	createdReview, createErr := controller.reviews.Create(
		ctx.UserContext(),
		order,
		body.Rating,
		body.Comment,
	)
	switch {
	case createErr == nil:
		return ctx.Status(http.StatusCreated).JSON(createdReview)
	case errors.Is(createErr, review.ErrOrderNotCompleted):
		return ctx.SendStatus(http.StatusPreconditionFailed)
	case errors.Is(createErr, review.ErrOrderAlreadyReviewed):
		return ctx.SendStatus(http.StatusConflict)
	default:
		return ctx.SendStatus(http.StatusInternalServerError)
	}

}

func (controller *Controller) Get(ctx *fiber.Ctx) error {
	order := ctx.Locals("order").(order.Order)

	rev, err := controller.reviews.GetByOrderID(ctx.UserContext(), order.ID)
	switch {
	case err == nil:
		return ctx.JSON(rev)
	case errors.Is(err, review.ErrNotFound):
		return ctx.SendStatus(http.StatusNotFound)
	default:
		return ctx.SendStatus(http.StatusInternalServerError)
	}
}

func (controller *Controller) Delete(ctx *fiber.Ctx) error {
	order := ctx.Locals("order").(order.Order)
	reviewToDelete, findErr := controller.reviews.GetByOrderID(ctx.UserContext(), order.ID)
	if findErr != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	deleteErr := controller.reviews.Delete(ctx.UserContext(), reviewToDelete)
	if deleteErr != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.SendStatus(http.StatusOK)
}

func (controller *Controller) middlewareGetOrder(ctx *fiber.Ctx) error {
	orderID := ctx.Params("orderID")
	if orderID == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	res, err := controller.orders.GetByID(ctx.UserContext(), orderID)
	switch {
	case err == nil:
		ctx.Locals("order", res)
		return ctx.Next()
	case errors.Is(err, order.ErrNotFound):
		return ctx.SendStatus(http.StatusNotFound)
	default:
		return ctx.SendStatus(http.StatusInternalServerError)
	}
}
