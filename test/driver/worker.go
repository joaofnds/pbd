package driver

import (
	"app/review"
	"app/test/matchers"
	"app/test/req"
	"app/worker"

	"net/http"
)

type WorkerDriver struct {
	url     string
	headers req.Headers
	Orders  *WorkerOrderDriver
}

func NewWorkerDriver(url string, headers req.Headers) *WorkerDriver {
	return &WorkerDriver{
		url:     url,
		headers: headers,
		Orders:  NewWorkerOrderDriver(url, headers),
	}
}

func (driver *WorkerDriver) Create(email, password string, hourlyRate int) (worker.Worker, error) {
	var cus worker.Worker

	return cus, makeJSONRequest(params{
		into:   &cus,
		status: http.StatusCreated,
		req: func() (*http.Response, error) {
			return req.Post(
				driver.url+"/workers",
				req.MergeHeaders(driver.headers, req.Headers{"Content-Type": "application/json"}),
				marshal(kv{"email": email, "password": password, "hourly_rate": hourlyRate}),
			)
		},
	})
}

func (driver *WorkerDriver) MustCreate(email, password string, hourlyRate int) worker.Worker {
	return matchers.Must2(driver.Create(email, password, hourlyRate))
}

func (driver *WorkerDriver) Get(id string) (worker.Worker, error) {
	var cus worker.Worker
	return cus, makeJSONRequest(params{
		into:   &cus,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				driver.url+"/workers/"+id,
				driver.headers,
			)
		},
	})
}

func (driver *WorkerDriver) MustGet(id string) worker.Worker {
	return matchers.Must2(driver.Get(id))
}

func (driver *WorkerDriver) ListReviews(id string) ([]review.Review, error) {
	var reviews []review.Review
	return reviews, makeJSONRequest(params{
		into:   &reviews,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				driver.url+"/workers/"+id+"/reviews",
				driver.headers,
			)
		},
	})
}

func (driver *WorkerDriver) MustListReviews(id string) []review.Review {
	return matchers.Must2(driver.ListReviews(id))
}
