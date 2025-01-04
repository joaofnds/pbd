package driver

import (
	"app/review"
	"app/review/review_http"
	"app/test/matchers"
	"app/test/req"

	"net/http"
)

type ReviewDriver struct {
	url     string
	headers req.Headers
}

func NewReviewDriver(url string, headers req.Headers) *ReviewDriver {
	return &ReviewDriver{
		url:     url,
		headers: headers,
	}
}

func (driver *ReviewDriver) Create(orderID string, createDTO review_http.CreateReviewDTO) (review.Review, error) {
	var cal review.Review

	return cal, makeJSONRequest(params{
		into:   &cal,
		status: http.StatusCreated,
		req: func() (*http.Response, error) {
			return req.Post(
				driver.url+"/orders/"+orderID+"/reviews",
				req.MergeHeaders(driver.headers, map[string]string{"Content-Type": "application/json"}),
				marshal(createDTO),
			)
		},
	})
}

func (driver *ReviewDriver) MustCreate(orderID string, createDTO review_http.CreateReviewDTO) review.Review {
	return matchers.Must2(driver.Create(orderID, createDTO))
}

func (driver *ReviewDriver) Get(orderID string) (review.Review, error) {
	var cal review.Review

	return cal, makeJSONRequest(params{
		into:   &cal,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				driver.url+"/orders/"+orderID+"/reviews",
				driver.headers,
			)
		},
	})
}

func (driver *ReviewDriver) MustGet(orderID string) review.Review {
	return matchers.Must2(driver.Get(orderID))
}

func (driver *ReviewDriver) Delete(orderID string) error {
	return makeJSONRequest(params{
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Delete(
				driver.url+"/orders/"+orderID+"/reviews",
				driver.headers,
			)
		},
	})
}

func (driver *ReviewDriver) MustDelete(orderID string) {
	matchers.Must(driver.Delete(orderID))
}
