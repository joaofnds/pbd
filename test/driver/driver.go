package driver

import (
	"app/test/matchers"
	"app/test/req"
	"app/user"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type Driver struct {
	url     string
	headers req.Headers

	Auth   *AuthDriver
	Health *HealthDriver
	KV     *KVDriver
	Users  *UserDriver

	Addresses *AddressDriver
	Calendars *CalendarDriver
	Customers *CustomerDriver
	Orders    *OrderDriver
	Phones    *PhoneDriver
	Workers   *WorkerDriver
}

func NewDriver(url string, headers req.Headers) *Driver {
	return &Driver{
		url:     url,
		headers: headers,
		Auth:    NewAuthDriver(url, headers),
		Health:  NewHealthDriver(url, headers),
		KV:      NewKVDriver(url, headers),
		Users:   NewUserDriver(url, headers),

		Addresses: NewAddressDriver(url, headers),
		Calendars: NewCalendarDriver(url, headers),
		Customers: NewCustomerDriver(url, headers),
		Orders:    NewOrderDriver(url, headers),
		Phones:    NewPhoneDriver(url, headers),
		Workers:   NewWorkerDriver(url, headers),
	}
}

func (driver *Driver) SetHeader(key, value string) {
	driver.headers[key] = value
}

func (driver *Driver) Login(email, password string) {
	token := driver.Auth.MustLogin(email, password)
	driver.headers.Set("Authorization", fmt.Sprintf("%s %s", token.TokenType, token.AccessToken))
}

type User struct {
	App    *Driver
	Entity user.User
}

type kv map[string]any

func marshal(v any) io.Reader {
	return bytes.NewReader(matchers.Must2(json.Marshal(v)))
}
