package driver

import (
	"app/calendar"
	"app/test/matchers"
	"app/test/req"

	"net/http"
)

type CalendarDriver struct {
	url     string
	headers req.Headers

	Events *CalendarEventsDriver
}

func NewCalendarDriver(url string, headers req.Headers) *CalendarDriver {
	return &CalendarDriver{
		url:     url,
		headers: headers,
		Events:  NewCalendarEventsDriver(url, headers),
	}
}

func (driver *CalendarDriver) Create(workerID string) (calendar.Calendar, error) {
	var cal calendar.Calendar

	return cal, makeJSONRequest(params{
		into:   &cal,
		status: http.StatusCreated,
		req: func() (*http.Response, error) {
			return req.Post(
				driver.url+"/workers/"+workerID+"/calendars",
				driver.headers,
				nil,
			)
		},
	})
}

func (driver *CalendarDriver) MustCreate(workerID string) calendar.Calendar {
	return matchers.Must2(driver.Create(workerID))
}

func (driver *CalendarDriver) Get(workerID string) (calendar.Calendar, error) {
	var cal calendar.Calendar

	return cal, makeJSONRequest(params{
		into:   &cal,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				driver.url+"/workers/"+workerID+"/calendars/"+workerID,
				driver.headers,
			)
		},
	})
}

func (driver *CalendarDriver) MustGet(workerID string) calendar.Calendar {
	return matchers.Must2(driver.Get(workerID))
}

func (driver *CalendarDriver) Delete(workerID string) error {
	return makeJSONRequest(params{
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Delete(
				driver.url+"/workers/"+workerID+"/calendars/"+workerID,
				driver.headers,
			)
		},
	})
}

func (driver *CalendarDriver) MustDelete(workerID string) {
	matchers.Must(driver.Delete(workerID))
}
