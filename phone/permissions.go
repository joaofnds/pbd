package phone

import (
	"app/authz"
	"app/internal/ref"
	"app/user"
)

func NewRef(id string) ref.Ref {
	return ref.New("phone", id)
}

type PermissionService struct {
	permissions authz.PermissionManager
}

func NewPermissionService(permissions authz.PermissionManager) *PermissionService {
	return &PermissionService{permissions: permissions}
}

func (service *PermissionService) GrantNewPhonePermissions(phone Phone) error {
	return service.permissions.Add(
		authz.NewAllowPolicy(user.NewRef(phone.UserID), NewRef(phone.ID), user.PermPhoneRead),
		authz.NewAllowPolicy(user.NewRef(phone.UserID), NewRef(phone.ID), user.PermPhoneUpdate),
		authz.NewAllowPolicy(user.NewRef(phone.UserID), NewRef(phone.ID), user.PermPhoneDelete),
	)
}
