package driver

import (
	"app/phone"
	"app/phone/phone_http"
	"app/test/matchers"
	"app/test/req"

	"net/http"
)

type PhoneDriver struct {
	url     string
	headers req.Headers
}

func NewPhoneDriver(url string, headers req.Headers) *PhoneDriver {
	return &PhoneDriver{url: url, headers: headers}
}

func (driver *PhoneDriver) Create(userID string, body phone_http.CreateBody) (phone.Phone, error) {
	var addr phone.Phone
	return addr, makeJSONRequest(params{
		into:   &addr,
		status: http.StatusCreated,
		req: func() (*http.Response, error) {
			return req.Post(
				driver.url+"/users/"+userID+"/phones",
				req.MergeHeaders(driver.headers, req.Headers{"Content-Type": "application/json"}),
				marshal(body),
			)
		},
	})
}

func (driver *PhoneDriver) MustCreate(userID string, body phone_http.CreateBody) phone.Phone {
	return matchers.Must2(driver.Create(userID, body))
}

func (driver *PhoneDriver) Get(userID, phoneID string) (phone.Phone, error) {
	var addr phone.Phone
	return addr, makeJSONRequest(params{
		into:   &addr,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				driver.url+"/users/"+userID+"/phones/"+phoneID,
				driver.headers,
			)
		},
	})
}

func (driver *PhoneDriver) MustGet(userID, phoneID string) phone.Phone {
	return matchers.Must2(driver.Get(userID, phoneID))
}

func (driver *PhoneDriver) List(userID string) ([]phone.Phone, error) {
	var addrs []phone.Phone
	return addrs, makeJSONRequest(params{
		into:   &addrs,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				driver.url+"/users/"+userID+"/phones",
				driver.headers,
			)
		},
	})
}

func (driver *PhoneDriver) MustList(userID string) []phone.Phone {
	return matchers.Must2(driver.List(userID))
}

func (driver *PhoneDriver) Delete(userID, phoneID string) error {
	_, err := makeRequest(
		http.StatusNoContent,
		func() (*http.Response, error) {
			return req.Delete(
				driver.url+"/users/"+userID+"/phones/"+phoneID,
				driver.headers,
			)
		},
	)
	return err
}

func (driver *PhoneDriver) MustDelete(userID, phoneID string) {
	matchers.Must(driver.Delete(userID, phoneID))
}
