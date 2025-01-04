package driver

import (
	"app/order"
	"app/order/order_http"
	"app/test/matchers"
	"app/test/req"

	"net/http"
)

type OrderDriver struct {
	url     string
	headers req.Headers

	Reviews *ReviewDriver
}

func NewOrderDriver(url string, headers req.Headers) *OrderDriver {
	return &OrderDriver{
		url:     url,
		headers: headers,

		Reviews: NewReviewDriver(url, headers),
	}
}

func (driver *OrderDriver) Create(dto order_http.CreatePayload) (order.Order, error) {
	var order order.Order
	return order, makeJSONRequest(params{
		into:   &order,
		status: http.StatusCreated,
		req: func() (*http.Response, error) {
			return req.Post(
				driver.url+"/orders",
				req.MergeHeaders(driver.headers, map[string]string{"Content-Type": "application/json"}),
				marshal(dto),
			)
		},
	})
}

func (driver *OrderDriver) MustCreate(dto order_http.CreatePayload) order.Order {
	return matchers.Must2(driver.Create(dto))
}

func (driver *OrderDriver) List() ([]order.Order, error) {
	var orders []order.Order
	return orders, makeJSONRequest(params{
		into:   &orders,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				driver.url+"/orders",
				driver.headers,
			)
		},
	})
}

func (driver *OrderDriver) MustList() []order.Order {
	return matchers.Must2(driver.List())
}

func (driver *OrderDriver) Get(id string) (order.Order, error) {
	var order order.Order
	return order, makeJSONRequest(params{
		into:   &order,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				driver.url+"/orders/"+id,
				driver.headers,
			)
		},
	})
}

func (driver *OrderDriver) MustGet(id string) order.Order {
	return matchers.Must2(driver.Get(id))
}

func (driver *OrderDriver) Update(id string, dto order_http.UpdatePayload) error {
	_, err := makeRequest(
		http.StatusOK,
		func() (*http.Response, error) {
			return req.Patch(
				driver.url+"/orders/"+id,
				req.MergeHeaders(driver.headers, map[string]string{"Content-Type": "application/json"}),
				marshal(dto),
			)
		},
	)

	return err
}

func (driver *OrderDriver) MustUpdate(id string, dto order_http.UpdatePayload) {
	matchers.Must(driver.Update(id, dto))
}

func (driver *OrderDriver) Delete(id string) error {
	return makeJSONRequest(params{
		status: http.StatusNoContent,
		req: func() (*http.Response, error) {
			return req.Delete(
				driver.url+"/orders/"+id,
				driver.headers,
			)
		},
	})
}

func (driver *OrderDriver) MustDelete(id string) {
	matchers.Must(driver.Delete(id))
}
