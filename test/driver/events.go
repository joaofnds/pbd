package driver

import (
	"app/event"
	"app/event/event_http"
	"app/test/matchers"
	"app/test/req"
	"bytes"
	"encoding/json"
	"time"

	"net/http"
	"net/url"
)

type CalendarEventsDriver struct {
	url     string
	headers req.Headers
}

func NewCalendarEventsDriver(url string, headers req.Headers) *CalendarEventsDriver {
	return &CalendarEventsDriver{
		url:     url,
		headers: headers,
	}
}

func (driver *CalendarEventsDriver) Create(workerID string, calendarID string, dto event_http.CreateBody) (event.Event, error) {
	var evt event.Event

	return evt, makeJSONRequest(params{
		into:   &evt,
		status: http.StatusCreated,
		req: func() (*http.Response, error) {
			return req.Post(
				driver.url+"/workers/"+workerID+"/calendars/"+calendarID+"/events",
				req.MergeHeaders(driver.headers, map[string]string{"Content-Type": "appliction/json"}),
				bytes.NewReader(matchers.Must2(json.Marshal(dto))))
		},
	})
}

func (driver *CalendarEventsDriver) MustCreate(workerID string, calendarID string, dto event_http.CreateBody) event.Event {
	return matchers.Must2(driver.Create(workerID, calendarID, dto))
}

func (driver *CalendarEventsDriver) Get(workerID string, calendarID string, eventID string) (event.Event, error) {
	var evt event.Event

	return evt, makeJSONRequest(params{
		into:   &evt,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				driver.url+"/workers/"+workerID+"/calendars/"+calendarID+"/events"+"/"+eventID,
				driver.headers,
			)
		},
	})
}

func (driver *CalendarEventsDriver) List(workerID string, calendarID string, dto event_http.TimeQuery) ([]event.Event, error) {
	var events []event.Event

	reqURL := matchers.Must2(url.Parse(driver.url))
	reqURL.Path = "/workers/" + workerID + "/calendars/" + calendarID + "/events"
	reqURL.RawQuery = url.Values{
		"starts_at": {dto.StartsAt.Format(time.RFC3339)},
		"ends_at":   {dto.EndsAt.Format(time.RFC3339)},
	}.Encode()

	return events, makeJSONRequest(params{
		into:   &events,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				reqURL.String(),
				driver.headers,
			)
		},
	})
}

func (driver *CalendarEventsDriver) MustList(workerID string, calendarID string, dto event_http.TimeQuery) []event.Event {
	return matchers.Must2(driver.List(workerID, calendarID, dto))
}

func (driver *CalendarEventsDriver) MustGet(workerID string, calendarID string, eventID string) event.Event {
	return matchers.Must2(driver.Get(workerID, calendarID, eventID))
}

func (driver *CalendarEventsDriver) Update(workerID string, calendarID string, eventID string, dto event_http.UpdateBody) error {
	return makeJSONRequest(params{
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Patch(
				driver.url+"/workers/"+workerID+"/calendars/"+calendarID+"/events"+"/"+eventID,
				req.MergeHeaders(driver.headers, map[string]string{"Content-Type": "application/json"}),
				bytes.NewReader(matchers.Must2(json.Marshal(dto))))
		},
	})
}

func (driver *CalendarEventsDriver) MustUpdate(workerID string, calendarID string, eventID string, dto event_http.UpdateBody) {
	matchers.Must(driver.Update(workerID, calendarID, eventID, dto))
}

func (driver *CalendarEventsDriver) Delete(workerID string, calendarID string, eventID string) error {
	return makeJSONRequest(params{
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Delete(
				driver.url+"/workers/"+workerID+"/calendars/"+calendarID+"/events"+"/"+eventID,
				driver.headers,
			)
		},
	})
}

func (driver *CalendarEventsDriver) MustDelete(workerID string, calendarID string, eventID string) {
	matchers.Must(driver.Delete(workerID, calendarID, eventID))
}
