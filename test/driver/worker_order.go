package driver

import (
	"app/order"
	"app/test/matchers"
	"app/test/req"

	"net/http"
)

type WorkerOrderDriver struct {
	url     string
	headers req.Headers
}

func NewWorkerOrderDriver(url string, headers req.Headers) *WorkerOrderDriver {
	return &WorkerOrderDriver{url: url, headers: headers}
}

func (driver *WorkerOrderDriver) List(workerID string) ([]order.Order, error) {
	var orders []order.Order
	return orders, makeJSONRequest(params{
		into:   &orders,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				driver.url+"/workers/"+workerID+"/orders",
				driver.headers,
			)
		},
	})
}

func (driver *WorkerOrderDriver) MustList(workerID string) []order.Order {
	return matchers.Must2(driver.List(workerID))
}
