package customer

import (
	"app/authz"
	"app/internal/ref"
	"app/user"
)

const (
	PermCustomerRead = "customer:read"

	PermAddressCreate = "customer:address:create"
	PermAddressRead   = "customer:address:read"
	PermAddressList   = "customer:address:list"
	PermAddressUpdate = "customer:address:update"
	PermAddressDelete = "customer:address:delete"
)

func NewRef(id string) ref.Ref {
	return ref.New("customer", id)
}

type PermissionService struct {
	permissions authz.PermissionManager
}

func NewPermissionService(permissions authz.PermissionManager) *PermissionService {
	return &PermissionService{permissions: permissions}
}

func (service *PermissionService) GrantNewCustomerPermissions(customer Customer) error {
	return service.permissions.Add(
		authz.NewAllowPolicy(user.NewRef(customer.ID), NewRef(customer.ID), "*"),
	)
}
