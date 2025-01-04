package customer_http_test

import (
	"testing"
	"time"

	"app/event"
	"app/event/event_http"
	"app/order/order_http"
	"app/test/harness"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCustomerHTTP(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "customer http suite")
}

var _ = Describe("/customers", Ordered, func() {
	var app *harness.Harness
	t0 := time.Now().Truncate(time.Hour)

	BeforeAll(func() { app = harness.Setup() })
	BeforeEach(func() { app.BeforeEach() })
	AfterEach(func() { app.AfterEach() })
	AfterAll(func() { app.Teardown() })

	It("creates and gets customer", func() {
		api := app.NewDriver()
		customer := api.Customers.MustCreate("bob@example.com", "p455w0rd")
		api.Login("bob@example.com", "p455w0rd")

		Expect(api.Customers.MustGet(customer.ID)).To(Equal(customer))
	})

	Describe("orders", func() {
		It("lists orders", func() {
			worker := app.NewWorker("bob@example.com", "p455w0rd", 10_00)
			cal := worker.App.Calendars.MustGet(worker.Entity.ID)
			evt := worker.App.Calendars.Events.MustCreate(worker.Entity.ID, cal.ID, event_http.CreateBody{
				Status:   event.StatusAvailable,
				StartsAt: t0.Add(1 * time.Hour),
				EndsAt:   t0.Add(2 * time.Hour),
			})

			customer := app.NewCustomer("alice@example.com", "p455w0rd")
			addr := customer.App.Addresses.MustCreateTestAddress(customer.Entity.ID)
			ord := customer.App.Orders.MustCreate(order_http.CreatePayload{
				AddressID:  addr.ID,
				EventID:    evt.ID,
				WorkerID:   worker.Entity.ID,
				CustomerID: customer.Entity.ID,
				StartsAt:   evt.StartsAt,
				EndsAt:     evt.EndsAt,
			})

			customerOrders := customer.App.Customers.Orders.MustList(customer.Entity.ID)
			Expect(customerOrders).To(ConsistOf(ord))
		})
	})
})
