package driver

import (
	"app/customer"
	"app/review"
	"app/test/matchers"
	"app/test/req"

	"net/http"
)

type CustomerDriver struct {
	url     string
	headers req.Headers
	Orders  *CustomerOrderDriver
}

func NewCustomerDriver(url string, headers req.Headers) *CustomerDriver {
	return &CustomerDriver{
		url:     url,
		headers: headers,
		Orders:  NewCustomerOrderDriver(url, headers),
	}
}

func (driver *CustomerDriver) Create(email, password string) (customer.Customer, error) {
	var cus customer.Customer

	return cus, makeJSONRequest(params{
		into:   &cus,
		status: http.StatusCreated,
		req: func() (*http.Response, error) {
			return req.Post(
				driver.url+"/customers",
				req.MergeHeaders(driver.headers, req.Headers{"Content-Type": "application/json"}),
				marshal(kv{"email": email, "password": password}),
			)
		},
	})
}

func (driver *CustomerDriver) MustCreate(email, password string) customer.Customer {
	return matchers.Must2(driver.Create(email, password))
}

func (driver *CustomerDriver) Get(id string) (customer.Customer, error) {
	var cus customer.Customer
	return cus, makeJSONRequest(params{
		into:   &cus,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				driver.url+"/customers/"+id,
				driver.headers,
			)
		},
	})
}

func (driver *CustomerDriver) MustGet(id string) customer.Customer {
	return matchers.Must2(driver.Get(id))
}

func (driver *CustomerDriver) ListReviews(id string) ([]review.Review, error) {
	var reviews []review.Review
	return reviews, makeJSONRequest(params{
		into:   &reviews,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				driver.url+"/customers/"+id+"/reviews",
				driver.headers,
			)
		},
	})
}

func (driver *CustomerDriver) MustListReviews(id string) []review.Review {
	return matchers.Must2(driver.ListReviews(id))
}
