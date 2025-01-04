package event_http

import (
	"errors"
	"net/http"

	"app/adapter/validation"
	"app/calendar"
	"app/event"
	"app/worker"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func NewController(
	validator *validator.Validate,
	workers *worker.Service,
	calendars *calendar.Service,
	events *event.Service,
) *Controller {
	return &Controller{
		validator: validator,
		workers:   workers,
		calendars: calendars,
		events:    events,
	}
}

type Controller struct {
	validator *validator.Validate
	workers   *worker.Service
	calendars *calendar.Service
	events    *event.Service
}

func (controller *Controller) Register(app *fiber.App) {
	events := app.Group(
		"/workers/:workerID/calendars/:calendarID/events",
		controller.middlewareGetWorker,
		controller.middlewareGetCalendar,
	).Post("/", controller.Create).
		Get("/", controller.List)

	events.Group("/:eventID").
		Get("/", controller.Get).
		Patch("/", controller.Update).
		Delete("/", controller.Delete)
}

func (controller *Controller) Create(ctx *fiber.Ctx) error {
	var body CreateBody
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	if err := controller.validator.Struct(body); err != nil {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(fiber.Map{"errors": validation.ErrorMessages(err)})
	}

	cal := ctx.Locals("calendar").(calendar.Calendar)
	evt, err := controller.events.Create(ctx.UserContext(), event.CreateDTO{
		CalendarID: cal.ID,
		Status:     body.Status,
		StartsAt:   body.StartsAt,
		EndsAt:     body.EndsAt,
	})
	switch {
	case err == nil:
		return ctx.Status(http.StatusCreated).JSON(evt)
	case errors.Is(err, event.ErrTimeSlotTaken):
		return ctx.Status(http.StatusConflict).SendString(err.Error())
	default:
		return ctx.SendStatus(http.StatusInternalServerError)
	}
}

func (controller *Controller) List(ctx *fiber.Ctx) error {
	var params TimeQuery
	if err := ctx.QueryParser(&params); err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	if err := controller.validator.Struct(params); err != nil {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(fiber.Map{"errors": validation.ErrorMessages(err)})
	}

	cal := ctx.Locals("calendar").(calendar.Calendar)
	events, listErr := controller.events.List(ctx.UserContext(), event.ListDTO{
		CalendarID: cal.ID,
		StartsAt:   params.StartsAt,
		EndsAt:     params.EndsAt,
	})
	if listErr != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.JSON(events)
}

func (controller *Controller) Get(ctx *fiber.Ctx) error {
	eventID := ctx.Params("eventID")
	if eventID == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	evt, err := controller.events.FindByID(ctx.UserContext(), eventID)
	if err != nil {
		if errors.Is(err, event.ErrNotFound) {
			return ctx.SendStatus(http.StatusNotFound)
		} else {
			return ctx.SendStatus(http.StatusInternalServerError)
		}
	}

	return ctx.JSON(evt)
}

func (controller *Controller) Update(ctx *fiber.Ctx) error {
	var body UpdateBody
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	if err := controller.validator.Struct(body); err != nil {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(fiber.Map{"errors": validation.ErrorMessages(err)})
	}

	eventID := ctx.Params("eventID")
	evt, findEventErr := controller.events.FindByID(ctx.UserContext(), eventID)
	if findEventErr != nil {
		if errors.Is(findEventErr, event.ErrNotFound) {
			return ctx.SendStatus(http.StatusNotFound)
		} else {
			return ctx.SendStatus(http.StatusInternalServerError)
		}
	}

	updateErr := controller.events.Update(ctx.UserContext(), evt, event.UpdateDTO{
		Status: body.Status,
	})
	if updateErr != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.SendStatus(http.StatusOK)
}

func (controller *Controller) Delete(ctx *fiber.Ctx) error {
	eventID := ctx.Params("eventID")
	evt, findEventErr := controller.events.FindByID(ctx.UserContext(), eventID)
	if findEventErr != nil {
		if errors.Is(findEventErr, event.ErrNotFound) {
			return ctx.SendStatus(http.StatusNotFound)
		} else {
			return ctx.SendStatus(http.StatusInternalServerError)
		}
	}

	deleteErr := controller.events.Delete(ctx.UserContext(), evt)
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

func (controller *Controller) middlewareGetCalendar(ctx *fiber.Ctx) error {
	calendarID := ctx.Params("calendarID")
	if calendarID == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	cal, err := controller.calendars.FindByID(ctx.UserContext(), calendarID)
	if err != nil {
		if errors.Is(err, worker.ErrNotFound) {
			return ctx.SendStatus(http.StatusNotFound)
		} else {
			return ctx.SendStatus(http.StatusInternalServerError)
		}
	}

	ctx.Locals("calendar", cal)

	return ctx.Next()
}
