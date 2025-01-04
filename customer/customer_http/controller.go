package customer_http

import (
	"errors"
	"net/http"

	"app/adapter/validation"
	"app/authn/authn_http"
	"app/authz/authz_http"
	"app/customer"
	"app/order"
	"app/review"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func NewController(
	validator *validator.Validate,
	authn *authn_http.AuthMiddleware,
	authz *authz_http.Middleware,
	customers *customer.Service,
	orders *order.Service,
	reviews *review.Service,
	registration *customer.RegistrationService,
) *Controller {
	return &Controller{
		validator:    validator,
		authn:        authn,
		authz:        authz,
		customers:    customers,
		orders:       orders,
		reviews:      reviews,
		registration: registration,
	}
}

type Controller struct {
	validator    *validator.Validate
	authn        *authn_http.AuthMiddleware
	authz        *authz_http.Middleware
	customers    *customer.Service
	orders       *order.Service
	reviews      *review.Service
	registration *customer.RegistrationService
}

func (controller *Controller) Register(app *fiber.App) {
	customers := app.Group("/customers")
	customers.Post("/", controller.Create)

	customerRoute := app.Group(
		"/customers/:customerID",
		controller.authn.RequireUser,
		controller.middlewareGetCustomer,
	)
	customerRoute.Get(
		"/",
		controller.authz.RequireParamPermission("customer:customerID", customer.PermCustomerRead),
		controller.Get,
	)

	customerOrders := customerRoute.Group("/orders")
	customerOrders.Get(
		"/",
		controller.authz.RequireParamPermission("customer:customerID", order.PermOrderList),
		controller.ListOrders,
	)

	customerReviews := customerRoute.Group("/reviews")
	customerReviews.Get("/", controller.ListReviews)
}

func (controller *Controller) Create(ctx *fiber.Ctx) error {
	var body CreateCustomerDTO
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	if err := controller.validator.Struct(body); err != nil {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(fiber.Map{"errors": validation.ErrorMessages(err)})
	}

	cus, err := controller.registration.Register(ctx.UserContext(), body.Email, body.Password)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.Status(http.StatusCreated).JSON(cus)
}

func (controller *Controller) Get(ctx *fiber.Ctx) error {
	return ctx.JSON(ctx.Locals("customer"))
}

func (controller *Controller) ListOrders(ctx *fiber.Ctx) error {
	customer := ctx.Locals("customer").(customer.Customer)

	orders, err := controller.orders.FindByCustomerID(ctx.UserContext(), customer.ID)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.JSON(orders)
}

func (controller *Controller) ListReviews(ctx *fiber.Ctx) error {
	customer := ctx.Locals("customer").(customer.Customer)

	reviews, err := controller.reviews.FindByCustomerID(ctx.UserContext(), customer.ID)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.JSON(reviews)
}

func (controller *Controller) middlewareGetCustomer(ctx *fiber.Ctx) error {
	customerID := ctx.Params("customerID")
	if customerID == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	customerFound, err := controller.customers.FindByID(ctx.UserContext(), customerID)
	if err != nil {
		if errors.Is(err, customer.ErrNotFound) {
			return ctx.SendStatus(http.StatusNotFound)
		} else {
			return ctx.SendStatus(http.StatusInternalServerError)
		}
	}

	ctx.Locals("customer", customerFound)
	return ctx.Next()
}
