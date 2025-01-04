package order_http_test

import (
	"app/event"
	"app/event/event_http"
	"app/order"
	"app/order/order_http"
	"net/http"

	"app/test/driver"
	"app/test/harness"
	"testing"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestOrderHTTP(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "order http suite")
}

var _ = Describe("/order", Ordered, func() {
	var app *harness.Harness
	t0 := time.Now().Truncate(time.Hour).Add(1 * time.Hour)

	BeforeAll(func() { app = harness.Setup() })
	BeforeEach(func() { app.BeforeEach() })
	AfterEach(func() { app.AfterEach() })
	AfterAll(func() { app.Teardown() })

	It("books on top of an available event", func() {
		worker := app.NewWorker("bob@example.com", "p455w0rd", 10_00)
		cal := worker.App.Calendars.MustGet(worker.Entity.ID)
		evt := worker.App.Calendars.Events.MustCreate(worker.Entity.ID, cal.ID, event_http.CreateBody{
			Status:   event.StatusAvailable,
			StartsAt: t0.Add(1 * time.Hour),
			EndsAt:   t0.Add(3 * time.Hour),
		})

		customer := app.NewCustomer("alice@example.com", "p455w0rd")
		addr := customer.App.Addresses.MustCreateTestAddress(customer.Entity.ID)
		ord, err := customer.App.Orders.Create(order_http.CreatePayload{
			AddressID:  addr.ID,
			EventID:    evt.ID,
			WorkerID:   worker.Entity.ID,
			CustomerID: customer.Entity.ID,
			StartsAt:   evt.StartsAt,
			EndsAt:     evt.EndsAt,
		})

		Expect(err).To(BeNil())
		Expect(ord.ID).To(HaveLen(36))
		Expect(ord.Price).To(Equal(20_00)) // 2 hours at $10/hour
		Expect(ord.Status).To(Equal(order.StatusCreated))
		Expect(ord.EventID).NotTo(Equal(evt.ID))
		Expect(ord.WorkerID).To(Equal(worker.Entity.ID))
		Expect(ord.CustomerID).To(Equal(customer.Entity.ID))
		Expect(ord.CreatedAt).To(BeTemporally("~", time.Now(), time.Second))
		Expect(ord.UpdatedAt).To(BeTemporally("~", time.Now(), time.Second))
	})

	It("updates the event", func() {
		worker := app.NewWorker("bob@example.com", "p455w0rd", 10_00)
		cal := worker.App.Calendars.MustGet(worker.Entity.ID)
		evt := worker.App.Calendars.Events.MustCreate(worker.Entity.ID, cal.ID, event_http.CreateBody{
			Status:   event.StatusAvailable,
			StartsAt: t0.Add(1 * time.Hour),
			EndsAt:   t0.Add(2 * time.Hour),
		})

		customer := app.NewCustomer("alice@example.com", "p455w0rd")
		addr := customer.App.Addresses.MustCreateTestAddress(customer.Entity.ID)
		customer.App.Orders.MustCreate(order_http.CreatePayload{
			AddressID:  addr.ID,
			EventID:    evt.ID,
			WorkerID:   worker.Entity.ID,
			CustomerID: customer.Entity.ID,
			StartsAt:   evt.StartsAt,
			EndsAt:     evt.EndsAt,
		})

		events := worker.App.Calendars.Events.MustList(worker.Entity.ID, cal.ID, event_http.TimeQuery{
			StartsAt: evt.StartsAt,
			EndsAt:   evt.EndsAt,
		})
		Expect(events).To(HaveLen(1))
		Expect(events[0].Status).To(Equal(event.StatusBooked))
		Expect(events[0].StartsAt).To(Equal(evt.StartsAt))
		Expect(events[0].EndsAt).To(Equal(evt.EndsAt))
	})

	When("address is missing", func() {
		It("returns an error", func() {
			worker := app.NewWorker("bob@example.com", "p455w0rd", 10_00)
			cal := worker.App.Calendars.MustGet(worker.Entity.ID)
			evt := worker.App.Calendars.Events.MustCreate(worker.Entity.ID, cal.ID, event_http.CreateBody{
				Status:   event.StatusAvailable,
				StartsAt: t0.Add(1 * time.Hour),
				EndsAt:   t0.Add(2 * time.Hour),
			})

			customer := app.NewCustomer("alice@example.com", "p455w0rd")
			addr := customer.App.Addresses.MustCreateTestAddress(customer.Entity.ID)
			customer.App.Orders.MustCreate(order_http.CreatePayload{
				AddressID:  addr.ID,
				EventID:    evt.ID,
				WorkerID:   worker.Entity.ID,
				CustomerID: customer.Entity.ID,
				StartsAt:   evt.StartsAt,
				EndsAt:     evt.EndsAt,
			})

			_, err := customer.App.Orders.Create(order_http.CreatePayload{
				AddressID:  uuid.NewString(),
				EventID:    evt.ID,
				WorkerID:   worker.Entity.ID,
				CustomerID: customer.Entity.ID,
				StartsAt:   evt.StartsAt,
				EndsAt:     evt.EndsAt,
			})

			Expect(err).To(MatchError(driver.RequestFailure{
				Status: http.StatusBadRequest,
				Body:   `{"error":"order error: address error: address not found"}`,
			}))
		})
	})

	When("ordering on top of a booked event", func() {
		It("returns an error", func() {
			worker := app.NewWorker("bob@example.com", "p455w0rd", 10_00)
			cal := worker.App.Calendars.MustGet(worker.Entity.ID)
			evt := worker.App.Calendars.Events.MustCreate(worker.Entity.ID, cal.ID, event_http.CreateBody{
				Status:   event.StatusBooked,
				StartsAt: t0.Add(1 * time.Hour),
				EndsAt:   t0.Add(2 * time.Hour),
			})

			customer := app.NewCustomer("alice@example.com", "p455w0rd")
			addr := customer.App.Addresses.MustCreateTestAddress(customer.Entity.ID)

			_, err := customer.App.Orders.Create(order_http.CreatePayload{
				AddressID:  addr.ID,
				EventID:    evt.ID,
				WorkerID:   worker.Entity.ID,
				CustomerID: customer.Entity.ID,
				StartsAt:   evt.StartsAt,
				EndsAt:     evt.EndsAt,
			})

			Expect(err).To(MatchError(driver.RequestFailure{
				Status: http.StatusBadRequest,
				Body:   `{"error":"order error: booking error: event is not available"}`,
			}))
		})
	})

	When("start is after end", func() {
		It("returns an error", func() {
			worker := app.NewWorker("bob@example.com", "p455w0rd", 10_00)
			cal := worker.App.Calendars.MustGet(worker.Entity.ID)
			evt := worker.App.Calendars.Events.MustCreate(worker.Entity.ID, cal.ID, event_http.CreateBody{
				Status:   event.StatusAvailable,
				StartsAt: t0.Add(1 * time.Hour),
				EndsAt:   t0.Add(2 * time.Hour),
			})

			customer := app.NewCustomer("alice@example.com", "p455w0rd")
			addr := customer.App.Addresses.MustCreateTestAddress(customer.Entity.ID)

			_, err := customer.App.Orders.Create(order_http.CreatePayload{
				AddressID:  addr.ID,
				EventID:    evt.ID,
				WorkerID:   worker.Entity.ID,
				CustomerID: customer.Entity.ID,
				StartsAt:   evt.EndsAt,
				EndsAt:     evt.StartsAt,
			})

			Expect(err).To(MatchError(driver.RequestFailure{
				Status: http.StatusBadRequest,
				Body:   `{"error":"order error: booking error: starts at must be before ends at"}`,
			}))
		})
	})

	When("start is equal to end", func() {
		It("returns an error", func() {
			worker := app.NewWorker("bob@example.com", "p455w0rd", 10_00)
			cal := worker.App.Calendars.MustGet(worker.Entity.ID)
			evt := worker.App.Calendars.Events.MustCreate(worker.Entity.ID, cal.ID, event_http.CreateBody{
				Status:   event.StatusAvailable,
				StartsAt: t0.Add(1 * time.Hour),
				EndsAt:   t0.Add(2 * time.Hour),
			})

			customer := app.NewCustomer("alice@example.com", "p455w0rd")
			addr := customer.App.Addresses.MustCreateTestAddress(customer.Entity.ID)

			_, err := customer.App.Orders.Create(order_http.CreatePayload{
				AddressID:  addr.ID,
				EventID:    evt.ID,
				WorkerID:   worker.Entity.ID,
				CustomerID: customer.Entity.ID,
				StartsAt:   evt.StartsAt,
				EndsAt:     evt.StartsAt,
			})

			Expect(err).To(MatchError(driver.RequestFailure{
				Status: http.StatusBadRequest,
				Body:   `{"error":"order error: booking error: starts at must be before ends at"}`,
			}))
		})
	})

	When("start is in the past", func() {
		It("returns an error", func() {
			worker := app.NewWorker("bob@example.com", "p455w0rd", 10_00)
			cal := worker.App.Calendars.MustGet(worker.Entity.ID)
			evt := worker.App.Calendars.Events.MustCreate(worker.Entity.ID, cal.ID, event_http.CreateBody{
				Status:   event.StatusAvailable,
				StartsAt: t0.Add(1 * time.Hour),
				EndsAt:   t0.Add(2 * time.Hour),
			})

			customer := app.NewCustomer("alice@example.com", "p455w0rd")
			addr := customer.App.Addresses.MustCreateTestAddress(customer.Entity.ID)

			_, err := customer.App.Orders.Create(order_http.CreatePayload{
				AddressID:  addr.ID,
				EventID:    evt.ID,
				WorkerID:   worker.Entity.ID,
				CustomerID: customer.Entity.ID,
				StartsAt:   evt.StartsAt.Add(-2 * time.Hour),
				EndsAt:     evt.EndsAt,
			})

			Expect(err).To(MatchError(driver.RequestFailure{
				Status: http.StatusBadRequest,
				Body:   `{"error":"order error: booking error: starts at must be in the future"}`,
			}))
		})
	})

	When("start is before event start", func() {
		It("returns an error", func() {
			worker := app.NewWorker("bob@example.com", "p455w0rd", 10_00)
			cal := worker.App.Calendars.MustGet(worker.Entity.ID)
			evt := worker.App.Calendars.Events.MustCreate(worker.Entity.ID, cal.ID, event_http.CreateBody{
				Status:   event.StatusAvailable,
				StartsAt: t0.Add(1 * time.Hour),
				EndsAt:   t0.Add(2 * time.Hour),
			})

			customer := app.NewCustomer("alice@example.com", "p455w0rd")
			addr := customer.App.Addresses.MustCreateTestAddress(customer.Entity.ID)

			_, err := customer.App.Orders.Create(order_http.CreatePayload{
				AddressID:  addr.ID,
				EventID:    evt.ID,
				WorkerID:   worker.Entity.ID,
				CustomerID: customer.Entity.ID,
				StartsAt:   evt.StartsAt.Add(-5 * time.Minute),
				EndsAt:     evt.EndsAt,
			})

			Expect(err).To(MatchError(driver.RequestFailure{
				Status: http.StatusBadRequest,
				Body:   `{"error":"order error: booking error: request is not within event time"}`,
			}))
		})
	})

	When("end is after event end", func() {
		It("returns an error", func() {
			worker := app.NewWorker("bob@example.com", "p455w0rd", 10_00)
			cal := worker.App.Calendars.MustGet(worker.Entity.ID)
			evt := worker.App.Calendars.Events.MustCreate(worker.Entity.ID, cal.ID, event_http.CreateBody{
				Status:   event.StatusAvailable,
				StartsAt: t0.Add(1 * time.Hour),
				EndsAt:   t0.Add(2 * time.Hour),
			})

			customer := app.NewCustomer("alice@example.com", "p455w0rd")
			addr := customer.App.Addresses.MustCreateTestAddress(customer.Entity.ID)

			_, err := customer.App.Orders.Create(order_http.CreatePayload{
				AddressID:  addr.ID,
				EventID:    evt.ID,
				WorkerID:   worker.Entity.ID,
				CustomerID: customer.Entity.ID,
				StartsAt:   evt.StartsAt,
				EndsAt:     evt.EndsAt.Add(5 * time.Minute),
			})

			Expect(err).To(MatchError(driver.RequestFailure{
				Status: http.StatusBadRequest,
				Body:   `{"error":"order error: booking error: request is not within event time"}`,
			}))
		})
	})

	invalidDurations := []time.Duration{
		30 * time.Minute,
		45 * time.Minute,
		1*time.Hour + 30*time.Minute,
	}
	for _, duration := range invalidDurations {
		When("duration is "+duration.String(), func() {
			It("returns an error", func() {
				worker := app.NewWorker("bob@example.com", "p455w0rd", 10_00)
				cal := worker.App.Calendars.MustGet(worker.Entity.ID)
				evt := worker.App.Calendars.Events.MustCreate(worker.Entity.ID, cal.ID, event_http.CreateBody{
					Status:   event.StatusAvailable,
					StartsAt: t0.Add(1 * time.Hour),
					EndsAt:   t0.Add(10 * time.Hour),
				})

				customer := app.NewCustomer("alice@example.com", "p455w0rd")
				addr := customer.App.Addresses.MustCreateTestAddress(customer.Entity.ID)

				_, err := customer.App.Orders.Create(order_http.CreatePayload{
					AddressID:  addr.ID,
					EventID:    evt.ID,
					WorkerID:   worker.Entity.ID,
					CustomerID: customer.Entity.ID,
					StartsAt:   evt.StartsAt,
					EndsAt:     evt.StartsAt.Add(duration),
				})

				Expect(err).To(MatchError(driver.RequestFailure{
					Status: http.StatusBadRequest,
					Body:   `{"error":"order error: booking error: invalid duration"}`,
				}))
			})
		})
	}

	It("lists orders", func() {
		worker := app.NewWorker("bob@example.com", "p455w0rd", 10_00)
		cal := worker.App.Calendars.MustGet(worker.Entity.ID)
		firstEvent := worker.App.Calendars.Events.MustCreate(worker.Entity.ID, cal.ID, event_http.CreateBody{
			Status:   event.StatusAvailable,
			StartsAt: t0.Add(1 * time.Hour),
			EndsAt:   t0.Add(2 * time.Hour),
		})
		secondEvent := worker.App.Calendars.Events.MustCreate(worker.Entity.ID, cal.ID, event_http.CreateBody{
			Status:   event.StatusAvailable,
			StartsAt: t0.Add(3 * time.Hour),
			EndsAt:   t0.Add(4 * time.Hour),
		})

		customer := app.NewCustomer("alice@example.com", "p455w0rd")
		addr := customer.App.Addresses.MustCreateTestAddress(customer.Entity.ID)

		Expect(customer.App.Orders.MustList()).To(HaveLen(0))

		firstOrder := customer.App.Orders.MustCreate(order_http.CreatePayload{
			AddressID:  addr.ID,
			EventID:    firstEvent.ID,
			WorkerID:   worker.Entity.ID,
			CustomerID: customer.Entity.ID,
			StartsAt:   firstEvent.StartsAt,
			EndsAt:     firstEvent.EndsAt,
		})
		secondOrder := customer.App.Orders.MustCreate(order_http.CreatePayload{
			AddressID:  addr.ID,
			EventID:    secondEvent.ID,
			WorkerID:   worker.Entity.ID,
			CustomerID: customer.Entity.ID,
			StartsAt:   secondEvent.StartsAt,
			EndsAt:     secondEvent.EndsAt,
		})

		Expect(customer.App.Orders.MustList()).To(Equal([]order.Order{firstOrder, secondOrder}))
	})

	It("finds order by ID", func() {
		worker := app.NewWorker("bob@example.com", "p455w0rd", 10_00)
		cal := worker.App.Calendars.MustGet(worker.Entity.ID)
		firstEvent := worker.App.Calendars.Events.MustCreate(worker.Entity.ID, cal.ID, event_http.CreateBody{
			Status:   event.StatusAvailable,
			StartsAt: t0.Add(1 * time.Hour),
			EndsAt:   t0.Add(2 * time.Hour),
		})
		secondEvent := worker.App.Calendars.Events.MustCreate(worker.Entity.ID, cal.ID, event_http.CreateBody{
			Status:   event.StatusAvailable,
			StartsAt: t0.Add(3 * time.Hour),
			EndsAt:   t0.Add(4 * time.Hour),
		})

		customer := app.NewCustomer("alice@example.com", "p455w0rd")
		addr := customer.App.Addresses.MustCreateTestAddress(customer.Entity.ID)

		firstOrder := customer.App.Orders.MustCreate(order_http.CreatePayload{
			AddressID:  addr.ID,
			EventID:    firstEvent.ID,
			WorkerID:   worker.Entity.ID,
			CustomerID: customer.Entity.ID,
			StartsAt:   firstEvent.StartsAt,
			EndsAt:     firstEvent.EndsAt,
		})
		secondOrder := customer.App.Orders.MustCreate(order_http.CreatePayload{
			AddressID:  addr.ID,
			EventID:    secondEvent.ID,
			WorkerID:   worker.Entity.ID,
			CustomerID: customer.Entity.ID,
			StartsAt:   secondEvent.StartsAt,
			EndsAt:     secondEvent.EndsAt,
		})

		Expect(customer.App.Orders.MustGet(firstOrder.ID)).To(Equal(firstOrder))
		Expect(customer.App.Orders.MustGet(secondOrder.ID)).To(Equal(secondOrder))
	})

	It("deletes order by ID", func() {
		worker := app.NewWorker("bob@example.com", "p455w0rd", 10_00)
		cal := worker.App.Calendars.MustGet(worker.Entity.ID)
		firstEvent := worker.App.Calendars.Events.MustCreate(worker.Entity.ID, cal.ID, event_http.CreateBody{
			Status:   event.StatusAvailable,
			StartsAt: t0.Add(1 * time.Hour),
			EndsAt:   t0.Add(2 * time.Hour),
		})
		secondEvent := worker.App.Calendars.Events.MustCreate(worker.Entity.ID, cal.ID, event_http.CreateBody{
			Status:   event.StatusAvailable,
			StartsAt: t0.Add(3 * time.Hour),
			EndsAt:   t0.Add(4 * time.Hour),
		})

		customer := app.NewCustomer("alice@example.com", "p455w0rd")
		addr := customer.App.Addresses.MustCreateTestAddress(customer.Entity.ID)

		firstOrder := customer.App.Orders.MustCreate(order_http.CreatePayload{
			AddressID:  addr.ID,
			EventID:    firstEvent.ID,
			WorkerID:   worker.Entity.ID,
			CustomerID: customer.Entity.ID,
			StartsAt:   firstEvent.StartsAt,
			EndsAt:     firstEvent.EndsAt,
		})
		secondOrder := customer.App.Orders.MustCreate(order_http.CreatePayload{
			AddressID:  addr.ID,
			EventID:    secondEvent.ID,
			WorkerID:   worker.Entity.ID,
			CustomerID: customer.Entity.ID,
			StartsAt:   secondEvent.StartsAt,
			EndsAt:     secondEvent.EndsAt,
		})

		Expect(customer.App.Orders.MustList()).To(ConsistOf(firstOrder, secondOrder))

		customer.App.Orders.MustDelete(firstOrder.ID)

		Expect(customer.App.Orders.MustList()).To(ConsistOf(secondOrder))
	})

	It("updates order status", func() {
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
		Expect(ord.Status).To(Equal(order.StatusCreated))

		customer.App.Orders.MustUpdate(ord.ID, order_http.UpdatePayload{
			Status: order.StatusCompleted,
		})

		ord = customer.App.Orders.MustGet(ord.ID)
		Expect(ord.Status).To(Equal(order.StatusCompleted))
	})
})
