package address

import (
	"app/authz"
	"app/customer"
	"app/internal/ref"
	"app/user"
)

func NewRef(id string) ref.Ref {
	return ref.New("address", id)
}

type PermissionService struct {
	permissions authz.PermissionManager
}

func NewPermissionService(permissions authz.PermissionManager) *PermissionService {
	return &PermissionService{permissions: permissions}
}

func (service *PermissionService) GrantNewAddressPermissions(address Address) error {
	return service.permissions.Add(
		authz.NewAllowPolicy(user.NewRef(address.CustomerID), NewRef(address.ID), customer.PermAddressRead),
		authz.NewAllowPolicy(user.NewRef(address.CustomerID), NewRef(address.ID), customer.PermAddressUpdate),
		authz.NewAllowPolicy(user.NewRef(address.CustomerID), NewRef(address.ID), customer.PermAddressDelete),
	)
}
