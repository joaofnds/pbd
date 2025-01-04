package event_http_test

import (
	"net/http"
	"testing"
	"time"

	"app/event"
	"app/event/event_http"
	"app/test/driver"
	"app/test/harness"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCalendarEventHTTP(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "calendar event http suite")
}

var _ = Describe("/workers/:workerID/calendars/:calendarID", Ordered, func() {
	var app *harness.Harness
	t0 := time.Now().Truncate(time.Second)

	BeforeAll(func() { app = harness.Setup() })
	BeforeEach(func() { app.BeforeEach() })
	AfterEach(func() { app.AfterEach() })
	AfterAll(func() { app.Teardown() })

	It("creates an event", func() {
		worker := app.NewWorker("bob@example.com", "p455w0rd", 10_00)
		cal := worker.App.Calendars.MustGet(worker.Entity.ID)
		dto := event_http.CreateBody{
			Status:   event.StatusAvailable,
			StartsAt: t0,
			EndsAt:   t0.Add(time.Hour),
		}

		evt := worker.App.Calendars.Events.MustCreate(worker.Entity.ID, cal.ID, dto)
		Expect(evt.ID).To(HaveLen(36))
		Expect(evt.CalendarID).To(Equal(cal.ID))
		Expect(evt.Status).To(Equal(dto.Status))
		Expect(evt.StartsAt).To(BeComparableTo(dto.StartsAt))
		Expect(evt.EndsAt).To(BeComparableTo(dto.EndsAt))
	})

	It("lists events", func() {
		worker := app.NewWorker("bob@example.com", "p455w0rd", 10_00)
		cal := worker.App.Calendars.MustGet(worker.Entity.ID)
		event1 := worker.App.Calendars.Events.MustCreate(worker.Entity.ID, cal.ID, event_http.CreateBody{
			Status:   event.StatusAvailable,
			StartsAt: t0,
			EndsAt:   t0.Add(time.Hour),
		})
		event2 := worker.App.Calendars.Events.MustCreate(worker.Entity.ID, cal.ID, event_http.CreateBody{
			Status:   event.StatusAvailable,
			StartsAt: t0.Add(2 * time.Hour),
			EndsAt:   t0.Add(3 * time.Hour),
		})

		events := worker.App.Calendars.Events.MustList(worker.Entity.ID, cal.ID, event_http.TimeQuery{
			StartsAt: event1.StartsAt,
			EndsAt:   event2.EndsAt,
		})
		Expect(events).To(ConsistOf(event1, event2))
	})

	It("finds an event", func() {
		worker := app.NewWorker("bob@example.com", "p455w0rd", 10_00)
		cal := worker.App.Calendars.MustGet(worker.Entity.ID)
		evt := worker.App.Calendars.Events.MustCreate(worker.Entity.ID, cal.ID, event_http.CreateBody{
			Status:   event.StatusAvailable,
			StartsAt: t0,
			EndsAt:   t0.Add(time.Hour),
		})

		foundEvent := worker.App.Calendars.Events.MustGet(worker.Entity.ID, cal.ID, evt.ID)
		Expect(foundEvent).To(Equal(evt))
	})

	It("updates an event", func() {
		worker := app.NewWorker("bob@example.com", "p455w0rd", 10_00)
		cal := worker.App.Calendars.MustGet(worker.Entity.ID)
		evt := worker.App.Calendars.Events.MustCreate(worker.Entity.ID, cal.ID, event_http.CreateBody{
			Status:   event.StatusAvailable,
			StartsAt: t0,
			EndsAt:   t0.Add(time.Hour),
		})

		worker.App.Calendars.Events.MustUpdate(worker.Entity.ID, cal.ID, evt.ID, event_http.UpdateBody{
			Status: event.StatusBooked,
		})

		foundEvent := worker.App.Calendars.Events.MustGet(worker.Entity.ID, cal.ID, evt.ID)
		Expect(foundEvent.Status).To(Equal(event.StatusBooked))
	})

	It("deletes an event", func() {
		worker := app.NewWorker("bob@example.com", "p455w0rd", 10_00)
		cal := worker.App.Calendars.MustGet(worker.Entity.ID)
		evt := worker.App.Calendars.Events.MustCreate(worker.Entity.ID, cal.ID, event_http.CreateBody{
			Status:   event.StatusAvailable,
			StartsAt: t0,
			EndsAt:   t0.Add(time.Hour),
		})

		worker.App.Calendars.Events.MustDelete(worker.Entity.ID, cal.ID, evt.ID)

		_, err := worker.App.Calendars.Events.Get(worker.Entity.ID, cal.ID, evt.ID)
		Expect(err).To(Equal(driver.RequestFailure{
			Status: http.StatusNotFound,
			Body:   "Not Found",
		}))
	})

	When("trying to create overlapping events", func() {
		nonOverlappableStatuses := []string{event.StatusAvailable, event.StatusBooked}

		for _, status := range nonOverlappableStatuses {
			When("event is "+status, func() {
				It("fails", func() {
					worker := app.NewWorker("bob@example.com", "p455w0rd", 10_00)
					cal := worker.App.Calendars.MustGet(worker.Entity.ID)
					evt := worker.App.Calendars.Events.MustCreate(worker.Entity.ID, cal.ID, event_http.CreateBody{
						Status:   status,
						StartsAt: t0,
						EndsAt:   t0.Add(time.Hour),
					})

					overlappingTimes := []struct{ StartsAt, EndsAt time.Time }{
						{evt.StartsAt, evt.EndsAt},                                    // same time
						{evt.StartsAt.Add(time.Minute), evt.EndsAt.Add(-time.Minute)}, // within
						{evt.StartsAt.Add(-time.Hour), evt.StartsAt.Add(time.Minute)}, // before with overlap
						{evt.EndsAt.Add(-time.Minute), evt.EndsAt.Add(time.Hour)},     // after with overlap
						{evt.StartsAt.Add(-time.Hour), evt.EndsAt.Add(time.Hour)},     // over the whole event
					}

					for _, times := range overlappingTimes {
						_, err := worker.App.Calendars.Events.Create(worker.Entity.ID, cal.ID, event_http.CreateBody{
							Status:   status,
							StartsAt: times.StartsAt,
							EndsAt:   times.EndsAt,
						})
						Expect(err).To(Equal(driver.RequestFailure{
							Status: http.StatusConflict,
							Body:   "time slot is taken",
						}))
					}

					Expect(worker.App.Calendars.Events.MustList(worker.Entity.ID, cal.ID, event_http.TimeQuery{
						StartsAt: evt.StartsAt.Add(-time.Hour),
						EndsAt:   evt.EndsAt.Add(time.Hour),
					})).To(ConsistOf(evt))
				})
			})
		}

		overlappableStatuses := []string{event.StatusCanceled}
		for _, status := range overlappableStatuses {
			When("event is "+status, func() {
				It("overlaps", func() {
					worker := app.NewWorker("bob@example.com", "p455w0rd", 10_00)
					cal := worker.App.Calendars.MustGet(worker.Entity.ID)
					evt := worker.App.Calendars.Events.MustCreate(worker.Entity.ID, cal.ID, event_http.CreateBody{
						Status:   status,
						StartsAt: t0,
						EndsAt:   t0.Add(time.Hour),
					})

					worker.App.Calendars.Events.MustCreate(worker.Entity.ID, cal.ID, event_http.CreateBody{
						Status:   status,
						StartsAt: evt.StartsAt,
						EndsAt:   evt.EndsAt,
					})

					Expect(worker.App.Calendars.Events.MustList(worker.Entity.ID, cal.ID, event_http.TimeQuery{
						StartsAt: evt.StartsAt.Add(-time.Hour),
						EndsAt:   evt.EndsAt.Add(time.Hour),
					})).To(HaveLen(2))
				})
			})
		}
	})
})
