package calendar_http_test

import (
	"net/http"
	"testing"

	"app/test/driver"
	"app/test/harness"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCalendarHTTP(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "calendar http suite")
}

var _ = Describe("/workers/:workerID/calendars", Ordered, func() {
	var app *harness.Harness

	BeforeAll(func() { app = harness.Setup() })
	BeforeEach(func() { app.BeforeEach() })
	AfterEach(func() { app.AfterEach() })
	AfterAll(func() { app.Teardown() })

	It("creates calendar", func() {
		worker := app.NewWorker("bob@example.com", "p455w0rd", 10_00)
		cal := worker.App.Calendars.MustGet(worker.Entity.ID)

		Expect(cal.ID).To(Equal(worker.Entity.ID))
	})

	It("gets calendar", func() {
		worker := app.NewWorker("bob@example.com", "p455w0rd", 10_00)
		cal := worker.App.Calendars.MustGet(worker.Entity.ID)

		Expect(worker.App.Calendars.MustGet(worker.Entity.ID)).To(Equal(cal))
	})

	It("deletes calendar", func() {
		worker := app.NewWorker("bob@example.com", "p455w0rd", 10_00)
		cal := worker.App.Calendars.MustGet(worker.Entity.ID)
		Expect(worker.App.Calendars.MustGet(worker.Entity.ID)).To(Equal(cal))

		worker.App.Calendars.MustDelete(worker.Entity.ID)

		_, err := worker.App.Calendars.Get(worker.Entity.ID)
		Expect(err).To(Equal(driver.RequestFailure{
			Status: http.StatusNotFound,
			Body:   "Not Found",
		}))
	})
})
