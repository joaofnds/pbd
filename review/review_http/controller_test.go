package review_http_test

import (
	"net/http"
	"testing"
	"time"

	"app/event"
	"app/event/event_http"
	"app/order"
	"app/order/order_http"
	"app/review/review_http"
	"app/test/driver"
	"app/test/harness"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestReviewHTTP(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "review http suite")
}

var _ = Describe("/orders/:orderID/reviews", Ordered, func() {
	var app *harness.Harness

	BeforeAll(func() { app = harness.Setup() })
	BeforeEach(func() { app.BeforeEach() })
	AfterEach(func() { app.AfterEach() })
	AfterAll(func() { app.Teardown() })

	createOrder := func(status string) (*driver.User, *driver.User, order.Order) {
		t0 := time.Now().Truncate(time.Hour)

		worker := app.NewWorker("bob@example.com", "p455w0rd", 10_00)
		cal := worker.App.Calendars.MustGet(worker.Entity.ID)
		evt := worker.App.Calendars.Events.MustCreate(worker.Entity.ID, cal.ID, event_http.CreateBody{
			Status:   event.StatusAvailable,
			StartsAt: t0.Add(1 * time.Hour),
			EndsAt:   t0.Add(2 * time.Hour),
		})

		customer := app.NewCustomer("alice@example.com", "p45ssw0rd")
		addr := customer.App.Addresses.MustCreateTestAddress(customer.Entity.ID)

		ord := customer.App.Orders.MustCreate(order_http.CreatePayload{
			AddressID:  addr.ID,
			EventID:    evt.ID,
			WorkerID:   worker.Entity.ID,
			CustomerID: customer.Entity.ID,
			StartsAt:   evt.StartsAt,
			EndsAt:     evt.EndsAt,
		})

		if status != order.StatusCreated {
			customer.App.Orders.MustUpdate(ord.ID, order_http.UpdatePayload{
				Status: status,
			})
		}

		return customer, worker, customer.App.Orders.MustGet(ord.ID)
	}

	When("order does not exist", func() {
		It("does not create the review", func() {
			customer := app.NewCustomer("alice@example.com", "p45ssw0rd")
			_, err := customer.App.Orders.Reviews.Create(uuid.NewString(), review_http.CreateReviewDTO{
				Rating:  50,
				Comment: "good",
			})

			Expect(err).To(Equal(driver.RequestFailure{
				Status: http.StatusNotFound,
				Body:   "Not Found",
			}))
		})
	})

	When("order is not completed", func() {
		It("does not create the review", func() {
			customer, _, ord := createOrder(order.StatusCreated)

			_, err := customer.App.Orders.Reviews.Create(ord.ID, review_http.CreateReviewDTO{
				Rating:  50,
				Comment: "good",
			})

			Expect(err).To(Equal(driver.RequestFailure{
				Status: http.StatusForbidden,
				Body:   "Forbidden",
			}))
		})
	})

	When("order is already reviewed", func() {
		It("does not create the review", func() {
			customer, _, ord := createOrder(order.StatusCompleted)

			customer.App.Orders.Reviews.MustCreate(ord.ID, review_http.CreateReviewDTO{
				Rating:  50,
				Comment: "good",
			})

			_, err := customer.App.Orders.Reviews.Create(ord.ID, review_http.CreateReviewDTO{
				Rating:  50,
				Comment: "good",
			})

			Expect(err).To(Equal(driver.RequestFailure{
				Status: http.StatusConflict,
				Body:   "Conflict",
			}))
		})
	})

	It("creates a review", func() {
		customer, _, ord := createOrder(order.StatusCompleted)

		rev, err := customer.App.Orders.Reviews.Create(ord.ID, review_http.CreateReviewDTO{
			Rating:  50,
			Comment: "good",
		})
		Expect(err).To(BeNil())

		Expect(rev.ID).To(HaveLen(36))
		Expect(rev.OrderID).To(Equal(ord.ID))
		Expect(rev.Rating).To(Equal(50))
		Expect(rev.Comment).To(Equal("good"))
		Expect(rev.CreatedAt).To(BeTemporally("~", time.Now(), time.Second))
		Expect(rev.UpdatedAt).To(BeTemporally("~", time.Now(), time.Second))
	})

	It("gets review for the order", func() {
		customer, worker, ord := createOrder(order.StatusCompleted)

		rev := customer.App.Orders.Reviews.MustCreate(ord.ID, review_http.CreateReviewDTO{
			Rating:  50,
			Comment: "good",
		})

		Expect(worker.App.Orders.Reviews.MustGet(ord.ID)).To(Equal(rev))
	})

	It("deletes review", func() {
		customer, worker, ord := createOrder(order.StatusCompleted)

		customer.App.Orders.Reviews.MustCreate(ord.ID, review_http.CreateReviewDTO{
			Rating:  50,
			Comment: "good",
		})

		customer.App.Orders.Reviews.MustDelete(ord.ID)

		_, err := worker.App.Orders.Reviews.Get(ord.ID)
		Expect(err).To(Equal(driver.RequestFailure{
			Status: http.StatusNotFound,
			Body:   "Not Found",
		}))
	})

	Describe("rating", func() {
		It("is required", func() {
			customer, _, ord := createOrder(order.StatusCompleted)

			_, err := customer.App.Orders.Reviews.Create(ord.ID, review_http.CreateReviewDTO{
				Comment: "good",
			})

			Expect(err).To(Equal(driver.RequestFailure{
				Status: http.StatusBadRequest,
				Body:   `{"errors":["Field validation for 'Rating' failed on the 'required' tag"]}`,
			}))
		})

		It("must be positive", func() {
			customer, _, ord := createOrder(order.StatusCompleted)

			_, err := customer.App.Orders.Reviews.Create(ord.ID, review_http.CreateReviewDTO{
				Rating:  -1,
				Comment: "good",
			})

			Expect(err).To(Equal(driver.RequestFailure{
				Status: http.StatusBadRequest,
				Body:   `{"errors":["Field validation for 'Rating' failed on the 'gte' tag"]}`,
			}))
		})

		It("must be less than or equal to 100", func() {
			customer, _, ord := createOrder(order.StatusCompleted)

			_, err := customer.App.Orders.Reviews.Create(ord.ID, review_http.CreateReviewDTO{
				Rating:  101,
				Comment: "good",
			})

			Expect(err).To(Equal(driver.RequestFailure{
				Status: http.StatusBadRequest,
				Body:   `{"errors":["Field validation for 'Rating' failed on the 'lte' tag"]}`,
			}))
		})
	})

	Describe("comment", func() {
		It("is optional", func() {
			customer, _, ord := createOrder(order.StatusCompleted)

			rev, err := customer.App.Orders.Reviews.Create(ord.ID, review_http.CreateReviewDTO{
				Rating: 50,
			})
			Expect(err).To(BeNil())

			Expect(rev.Comment).To(Equal(""))
		})
	})

	It("lists reviews for the customer", func() {
		customer, worker, ord := createOrder(order.StatusCompleted)

		rev := customer.App.Orders.Reviews.MustCreate(ord.ID, review_http.CreateReviewDTO{
			Rating:  50,
			Comment: "good",
		})

		Expect(customer.App.Customers.MustListReviews(ord.CustomerID)).To(ContainElement(rev))
		Expect(worker.App.Workers.MustListReviews(ord.WorkerID)).To(ContainElement(rev))
	})
})
