package phone_http_test

import (
	"testing"
	"time"

	"app/phone/phone_http"
	"app/test/driver"
	"app/test/harness"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPhoneHTTP(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "phone http suite")
}

var _ = Describe("/phones", Ordered, func() {
	var app *harness.Harness

	BeforeAll(func() { app = harness.Setup() })
	BeforeEach(func() { app.BeforeEach() })
	AfterEach(func() { app.AfterEach() })
	AfterAll(func() { app.Teardown() })

	It("creates and gets phone", func() {
		cus := app.NewCustomer("bob@example.com", "p455w0rd")
		createdPhone, err := cus.App.Phones.Create(cus.Entity.ID, phone_http.CreateBody{
			CountryCode: "55",
			AreaCode:    "11",
			Number:      "999999999",
		})
		Expect(err).To(BeNil())

		Expect(createdPhone.ID).To(HaveLen(36))
		Expect(createdPhone.CountryCode).To(Equal("55"))
		Expect(createdPhone.AreaCode).To(Equal("11"))
		Expect(createdPhone.Number).To(Equal("999999999"))
		Expect(createdPhone.CreatedAt).To(BeTemporally("~", time.Now(), 1*time.Second))
		Expect(createdPhone.UpdatedAt).To(BeTemporally("~", time.Now(), 1*time.Second))
	})

	It("gets phone", func() {
		cus := app.NewCustomer("bob@example.com", "p455w0rd")
		createPhone := cus.App.Phones.MustCreate(cus.Entity.ID, phone_http.CreateBody{
			CountryCode: "55",
			AreaCode:    "11",
			Number:      "999999999",
		})

		foundPhone, err := cus.App.Phones.Get(cus.Entity.ID, createPhone.ID)
		Expect(err).To(BeNil())
		Expect(foundPhone).To(Equal(createPhone))
	})

	It("deletes phone", func() {
		cus := app.NewCustomer("bob@example.com", "p455w0rd")
		createdPhone := cus.App.Phones.MustCreate(cus.Entity.ID, phone_http.CreateBody{
			CountryCode: "55",
			AreaCode:    "11",
			Number:      "999999999",
		})

		_, err := cus.App.Phones.Get(cus.Entity.ID, createdPhone.ID)
		Expect(err).To(BeNil())

		err = cus.App.Phones.Delete(cus.Entity.ID, createdPhone.ID)
		Expect(err).To(BeNil())

		_, err = cus.App.Phones.Get(cus.Entity.ID, createdPhone.ID)
		Expect(err).To(Equal(driver.RequestFailure{
			Status: 404,
			Body:   "Not Found",
		}))
	})

	It("lists phones", func() {
		cus := app.NewCustomer("bob@example.com", "p455w0rd")
		Expect(cus.App.Phones.MustList(cus.Entity.ID)).To(BeEmpty())

		phone1 := cus.App.Phones.MustCreate(cus.Entity.ID, phone_http.CreateBody{
			CountryCode: "55",
			AreaCode:    "11",
			Number:      "111111111",
		})
		Expect(cus.App.Phones.MustList(cus.Entity.ID)).To(ConsistOf(phone1))

		phone2 := cus.App.Phones.MustCreate(cus.Entity.ID, phone_http.CreateBody{
			CountryCode: "55",
			AreaCode:    "11",
			Number:      "222222222",
		})
		Expect(cus.App.Phones.MustList(cus.Entity.ID)).To(ConsistOf(phone1, phone2))
	})
})
