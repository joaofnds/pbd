package address_http

import (
	"errors"
	"fmt"
	"net/http"

	"app/adapter/validation"
	"app/address"
	"app/authn/authn_http"
	"app/authz/authz_http"
	"app/customer"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func NewController(
	validator *validator.Validate,
	authn *authn_http.AuthMiddleware,
	authz *authz_http.Middleware,
	customers *customer.Service,
	addresses *address.Service,
) *Controller {
	return &Controller{
		validator: validator,
		authn:     authn,
		authz:     authz,
		customers: customers,
		addresses: addresses,
	}
}

type Controller struct {
	validator *validator.Validate
	authn     *authn_http.AuthMiddleware
	authz     *authz_http.Middleware
	customers *customer.Service
	addresses *address.Service
}

func (controller *Controller) Register(app *fiber.App) {
	addresses := app.Group(
		"/customers/:customerID/addresses",
		controller.authn.RequireUser,
		controller.middlewareGetCustomer,
	)
	addresses.Post(
		"/",
		controller.authz.RequireParamPermission("customer:customerID", customer.PermAddressCreate),
		controller.createAddress,
	)
	addresses.Get(
		"/",
		controller.authz.RequireParamPermission("customer:customerID", customer.PermAddressList),
		controller.getAddresses,
	)

	addr := addresses.Group("/:addressID", controller.middlewareGetAddress)
	addr.Get(
		"/",
		controller.authz.RequireParamPermission("customer:customerID", customer.PermAddressCreate),
		controller.getAddress,
	)
	addr.Delete(
		"/",
		controller.authz.RequireParamPermission("customer:customerID", customer.PermAddressDelete),
		controller.deleteAddress,
	)
}

func (controller *Controller) createAddress(ctx *fiber.Ctx) error {
	var body CreateBody
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(fiber.Map{"errors": err.Error()})
	}

	if err := controller.validator.Struct(body); err != nil {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(fiber.Map{"errors": validation.ErrorMessages(err)})
	}

	customer := ctx.Locals("customer").(customer.Customer)
	newAddress, err := controller.addresses.Create(ctx.UserContext(), body.ToCreateDTO(customer.ID))
	switch {
	case err == nil:
		return ctx.Status(http.StatusCreated).JSON(newAddress)
	default:
		return ctx.SendStatus(http.StatusInternalServerError)
	}
}

func (controller *Controller) getAddress(ctx *fiber.Ctx) error {
	address := ctx.Locals("address").(address.Address)
	return ctx.JSON(address)
}

func (controller *Controller) getAddresses(ctx *fiber.Ctx) error {
	customer := ctx.Locals("customer").(customer.Customer)
	addresses, err := controller.addresses.FindByCustomerID(ctx.UserContext(), customer.ID)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}
	return ctx.JSON(addresses)
}

func (controller *Controller) deleteAddress(ctx *fiber.Ctx) error {
	address := ctx.Locals("address").(address.Address)
	err := controller.addresses.Delete(ctx.UserContext(), address.ID)
	if err != nil {
		fmt.Printf("%#v\n", err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}
	return ctx.SendStatus(http.StatusNoContent)
}

func (controller *Controller) middlewareGetCustomer(ctx *fiber.Ctx) error {
	customerID := ctx.Params("customerID")
	if customerID == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	customerFound, err := controller.customers.FindByID(ctx.UserContext(), customerID)
	switch {
	case err == nil:
		ctx.Locals("customer", customerFound)
		return ctx.Next()
	case errors.Is(err, customer.ErrNotFound):
		return ctx.SendStatus(http.StatusNotFound)
	default:
		return ctx.SendStatus(http.StatusInternalServerError)
	}
}

func (controller *Controller) middlewareGetAddress(ctx *fiber.Ctx) error {
	addressID := ctx.Params("addressID")
	if addressID == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	foundAddress, err := controller.addresses.Get(ctx.UserContext(), addressID)
	switch {
	case err == nil:
		ctx.Locals("address", foundAddress)
		return ctx.Next()
	case errors.Is(err, address.ErrNotFound):
		return ctx.SendStatus(http.StatusNotFound)
	default:
		return ctx.SendStatus(http.StatusInternalServerError)
	}
}
