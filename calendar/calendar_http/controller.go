package calendar_http

import (
	"errors"
	"net/http"

	"app/authn/authn_http"
	"app/authz/authz_http"
	"app/calendar"
	"app/worker"

	"github.com/gofiber/fiber/v2"
)

func NewController(
	authn *authn_http.AuthMiddleware,
	authz *authz_http.Middleware,
	workers *worker.Service,
	calendars *calendar.Service,
) *Controller {
	return &Controller{
		authn:     authn,
		authz:     authz,
		workers:   workers,
		calendars: calendars,
	}
}

type Controller struct {
	authn     *authn_http.AuthMiddleware
	authz     *authz_http.Middleware
	workers   *worker.Service
	calendars *calendar.Service
}

func (controller *Controller) Register(app *fiber.App) {
	calendars := app.Group(
		"/workers/:workerID/calendars",
		controller.middlewareGetWorker,
		controller.authn.RequireUser,
	)
	calendars.Post(
		"/",
		controller.authz.RequireParamPermission("worker:workerID", calendar.PermCalCreate),
		controller.Create,
	)

	calendars.Group("/:calendarID").
		Get(
			"/",
			controller.authz.RequireParamPermission("calendar:calendarID", calendar.PermCalRead),
			controller.Find,
		).
		Delete(
			"/",
			controller.authz.RequireParamPermission("calendar:calendarID", calendar.PermCalDelete),
			controller.Delete,
		)
}

func (controller *Controller) Create(ctx *fiber.Ctx) error {
	worker := ctx.Locals("worker").(worker.Worker)

	cal, err := controller.calendars.Create(ctx.UserContext(), worker.ID)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.Status(http.StatusCreated).JSON(cal)
}

func (controller *Controller) Find(ctx *fiber.Ctx) error {
	calendarID := ctx.Params("calendarID")
	if calendarID == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	calendars, err := controller.calendars.FindByID(ctx.UserContext(), calendarID)
	if err != nil {
		if errors.Is(err, calendar.ErrNotFound) {
			return ctx.SendStatus(http.StatusNotFound)
		}

		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.JSON(calendars)

}

func (controller *Controller) Delete(ctx *fiber.Ctx) error {
	calendarID := ctx.Params("calendarID")
	if calendarID == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	cal, findErr := controller.calendars.FindByID(ctx.UserContext(), calendarID)
	if findErr != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	deleteErr := controller.calendars.Delete(ctx.UserContext(), cal)
	if deleteErr != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.SendStatus(http.StatusOK)
}

func (controller *Controller) middlewareGetWorker(ctx *fiber.Ctx) error {
	workerID := ctx.Params("workerID")
	if workerID == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	res, err := controller.workers.FindByID(ctx.UserContext(), workerID)
	if err != nil {
		if errors.Is(err, worker.ErrNotFound) {
			return ctx.SendStatus(http.StatusNotFound)
		} else {
			return ctx.SendStatus(http.StatusInternalServerError)
		}
	}

	ctx.Locals("worker", res)
	return ctx.Next()
}
