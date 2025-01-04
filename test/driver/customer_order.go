package driver

import (
	"app/order"
	"app/test/matchers"
	"app/test/req"

	"net/http"
)

type CustomerOrderDriver struct {
	url     string
	headers req.Headers
}

func NewCustomerOrderDriver(url string, headers req.Headers) *CustomerOrderDriver {
	return &CustomerOrderDriver{url: url, headers: headers}
}

func (driver *CustomerOrderDriver) List(customerID string) ([]order.Order, error) {
	var orders []order.Order
	return orders, makeJSONRequest(params{
		into:   &orders,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				driver.url+"/customers/"+customerID+"/orders",
				driver.headers,
			)
		},
	})
}

func (driver *CustomerOrderDriver) MustList(customerID string) []order.Order {
	return matchers.Must2(driver.List(customerID))
}
