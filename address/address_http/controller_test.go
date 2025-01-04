package address_http_test

import (
	"testing"

	"app/address/address_http"
	"app/test/driver"
	"app/test/harness"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAddressHTTP(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "address http suite")
}

var _ = Describe("/addresss", Ordered, func() {
	var app *harness.Harness

	BeforeAll(func() { app = harness.Setup() })
	BeforeEach(func() { app.BeforeEach() })
	AfterEach(func() { app.AfterEach() })
	AfterAll(func() { app.Teardown() })

	It("creates and gets address", func() {
		cus := app.NewCustomer("bob@example.com", "p455w0rd")
		addr, err := cus.App.Addresses.Create(cus.Entity.ID, address_http.CreateBody{
			Street:       "Main St",
			Number:       "123",
			Complement:   "Apt 1",
			Neighborhood: "Downtown",
			City:         "Springfield",
			State:        "IL",
			ZipCode:      "62701",
			Country:      "USA",
		})
		Expect(err).To(BeNil())

		Expect(addr.ID).To(HaveLen(36))
		Expect(addr.CustomerID).To(Equal(cus.Entity.ID))
		Expect(addr.Street).To(Equal("Main St"))
		Expect(addr.Number).To(Equal("123"))
		Expect(addr.Complement).To(Equal("Apt 1"))
		Expect(addr.Neighborhood).To(Equal("Downtown"))
		Expect(addr.City).To(Equal("Springfield"))
		Expect(addr.State).To(Equal("IL"))
		Expect(addr.ZipCode).To(Equal("62701"))
		Expect(addr.Country).To(Equal("USA"))
	})

	It("gets address", func() {
		cus := app.NewCustomer("bob@example.com", "p455w0rd")
		addr := cus.App.Addresses.MustCreate(cus.Entity.ID, address_http.CreateBody{
			Street:       "Main St",
			Number:       "123",
			Complement:   "Apt 1",
			Neighborhood: "Downtown",
			City:         "Springfield",
			State:        "IL",
			ZipCode:      "62701",
			Country:      "USA",
		})

		foundAddr, err := cus.App.Addresses.Get(cus.Entity.ID, addr.ID)
		Expect(err).To(BeNil())
		Expect(foundAddr).To(Equal(addr))
	})

	It("deletes address", func() {
		cus := app.NewCustomer("bob@example.com", "p455w0rd")
		addr := cus.App.Addresses.MustCreate(cus.Entity.ID, address_http.CreateBody{
			Street:       "Main St",
			Number:       "123",
			Complement:   "Apt 1",
			Neighborhood: "Downtown",
			City:         "Springfield",
			State:        "IL",
			ZipCode:      "62701",
			Country:      "USA",
		})

		_, err := cus.App.Addresses.Get(cus.Entity.ID, addr.ID)
		Expect(err).To(BeNil())

		err = cus.App.Addresses.Delete(cus.Entity.ID, addr.ID)
		Expect(err).To(BeNil())

		_, err = cus.App.Addresses.Get(cus.Entity.ID, addr.ID)
		Expect(err).To(Equal(driver.RequestFailure{
			Status: 404,
			Body:   "Not Found",
		}))
	})

	It("lists addresses", func() {
		cus := app.NewCustomer("bob@example.com", "p455w0rd")
		Expect(cus.App.Addresses.MustList(cus.Entity.ID)).To(BeEmpty())

		addr1 := cus.App.Addresses.MustCreate(cus.Entity.ID, address_http.CreateBody{
			Street:       "Main St",
			Number:       "123",
			Complement:   "Apt 1",
			Neighborhood: "Downtown",
			City:         "Springfield",
			State:        "IL",
			ZipCode:      "62701",
			Country:      "USA",
		})
		Expect(cus.App.Addresses.MustList(cus.Entity.ID)).To(ConsistOf(addr1))

		addr2 := cus.App.Addresses.MustCreate(cus.Entity.ID, address_http.CreateBody{
			Street:       "Main St",
			Number:       "123",
			Complement:   "Apt 1",
			Neighborhood: "Downtown",
			City:         "Springfield",
			State:        "IL",
			ZipCode:      "62701",
			Country:      "USA",
		})
		Expect(cus.App.Addresses.MustList(cus.Entity.ID)).To(ConsistOf(addr1, addr2))
	})
})
