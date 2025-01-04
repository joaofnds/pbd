package phone_http

import (
	"errors"
	"net/http"

	"app/adapter/validation"
	"app/authn/authn_http"
	"app/authz/authz_http"
	"app/phone"
	"app/user"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func NewController(
	validator *validator.Validate,
	authn *authn_http.AuthMiddleware,
	authz *authz_http.Middleware,
	users *user.Service,
	phones *phone.Service,
) *Controller {
	return &Controller{
		validator: validator,
		authn:     authn,
		authz:     authz,
		users:     users,
		phones:    phones,
	}
}

type Controller struct {
	validator *validator.Validate
	authn     *authn_http.AuthMiddleware
	authz     *authz_http.Middleware
	users     *user.Service
	phones    *phone.Service
}

func (controller *Controller) Register(app *fiber.App) {
	phones := app.Group(
		"/users/:userID/phones",
		controller.authn.RequireUser,
		controller.middlewareGetUser,
	)
	phones.Post(
		"/",
		controller.authz.RequireParamPermission("user:userID", user.PermPhoneCreate),
		controller.createPhone,
	)
	phones.Get(
		"/",
		controller.authz.RequireParamPermission("user:userID", user.PermPhoneList),
		controller.getPhones,
	)

	phon := phones.Group("/:phoneID", controller.middlewareGetPhone)
	phon.Get(
		"/",
		controller.authz.RequireParamPermission("user:userID", user.PermPhoneCreate),
		controller.getPhone,
	)
	phon.Delete(
		"/",
		controller.authz.RequireParamPermission("user:userID", user.PermPhoneDelete),
		controller.deletePhone,
	)
}

func (controller *Controller) createPhone(ctx *fiber.Ctx) error {
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

	user := ctx.Locals("user").(user.User)
	newPhone, err := controller.phones.Create(ctx.UserContext(), body.ToCreateDTO(user.ID))
	switch {
	case err == nil:
		return ctx.Status(http.StatusCreated).JSON(newPhone)
	default:
		return ctx.SendStatus(http.StatusInternalServerError)
	}
}

func (controller *Controller) getPhone(ctx *fiber.Ctx) error {
	phone := ctx.Locals("phone").(phone.Phone)
	return ctx.JSON(phone)
}

func (controller *Controller) getPhones(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(user.User)
	phones, err := controller.phones.FindByUserID(ctx.UserContext(), user.ID)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}
	return ctx.JSON(phones)
}

func (controller *Controller) deletePhone(ctx *fiber.Ctx) error {
	phone := ctx.Locals("phone").(phone.Phone)
	err := controller.phones.Delete(ctx.UserContext(), phone.ID)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}
	return ctx.SendStatus(http.StatusNoContent)
}

func (controller *Controller) middlewareGetUser(ctx *fiber.Ctx) error {
	userID := ctx.Params("userID")
	if userID == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	userFound, err := controller.users.FindByID(ctx.UserContext(), userID)
	switch {
	case err == nil:
		ctx.Locals("user", userFound)
		return ctx.Next()
	case errors.Is(err, user.ErrNotFound):
		return ctx.SendStatus(http.StatusNotFound)
	default:
		return ctx.SendStatus(http.StatusInternalServerError)
	}
}

func (controller *Controller) middlewareGetPhone(ctx *fiber.Ctx) error {
	phoneID := ctx.Params("phoneID")
	if phoneID == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	foundPhone, err := controller.phones.Get(ctx.UserContext(), phoneID)
	switch {
	case err == nil:
		ctx.Locals("phone", foundPhone)
		return ctx.Next()
	case errors.Is(err, phone.ErrNotFound):
		return ctx.SendStatus(http.StatusNotFound)
	default:
		return ctx.SendStatus(http.StatusInternalServerError)
	}
}
