package driver

import (
	"app/address"
	"app/address/address_http"
	"app/test/matchers"
	"app/test/req"

	"net/http"
)

type AddressDriver struct {
	url     string
	headers req.Headers
}

func NewAddressDriver(url string, headers req.Headers) *AddressDriver {
	return &AddressDriver{url: url, headers: headers}
}

func (driver *AddressDriver) Create(customerID string, body address_http.CreateBody) (address.Address, error) {
	var addr address.Address
	return addr, makeJSONRequest(params{
		into:   &addr,
		status: http.StatusCreated,
		req: func() (*http.Response, error) {
			return req.Post(
				driver.url+"/customers/"+customerID+"/addresses",
				req.MergeHeaders(driver.headers, req.Headers{"Content-Type": "application/json"}),
				marshal(body),
			)
		},
	})
}

func (driver *AddressDriver) MustCreate(customerID string, body address_http.CreateBody) address.Address {
	return matchers.Must2(driver.Create(customerID, body))
}

func (driver *AddressDriver) MustCreateTestAddress(customerID string) address.Address {
	return driver.MustCreate(customerID, address_http.CreateBody{
		Street:       "Main St",
		Number:       "123",
		Complement:   "Apt 1",
		Neighborhood: "Downtown",
		City:         "Springfield",
		State:        "IL",
		ZipCode:      "62701",
		Country:      "USA",
	})
}

func (driver *AddressDriver) Get(customerID, addressID string) (address.Address, error) {
	var addr address.Address
	return addr, makeJSONRequest(params{
		into:   &addr,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				driver.url+"/customers/"+customerID+"/addresses/"+addressID,
				driver.headers,
			)
		},
	})
}

func (driver *AddressDriver) MustGet(customerID, addressID string) address.Address {
	return matchers.Must2(driver.Get(customerID, addressID))
}

func (driver *AddressDriver) List(customerID string) ([]address.Address, error) {
	var addrs []address.Address
	return addrs, makeJSONRequest(params{
		into:   &addrs,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				driver.url+"/customers/"+customerID+"/addresses",
				driver.headers,
			)
		},
	})
}

func (driver *AddressDriver) MustList(customerID string) []address.Address {
	return matchers.Must2(driver.List(customerID))
}

func (driver *AddressDriver) Delete(customerID, addressID string) error {
	_, err := makeRequest(
		http.StatusNoContent,
		func() (*http.Response, error) {
			return req.Delete(
				driver.url+"/customers/"+customerID+"/addresses/"+addressID,
				driver.headers,
			)
		},
	)
	return err
}

func (driver *AddressDriver) MustDelete(customerID, addressID string) {
	matchers.Must(driver.Delete(customerID, addressID))
}
