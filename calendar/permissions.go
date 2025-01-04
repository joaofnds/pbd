package calendar

import (
	"app/authz"
	"app/internal/ref"
	"app/user"
)

const (
	PermCalCreate = "calendar:create"
	PermCalRead   = "calendar:read"
	PermCalDelete = "calendar:delete"
)

func NewRef(id string) ref.Ref {
	return ref.New("calendar", id)
}

type PermissionService struct {
	permissions authz.PermissionManager
}

func NewPermissionService(permissions authz.PermissionManager) *PermissionService {
	return &PermissionService{permissions: permissions}
}

func (service *PermissionService) NewCalendarPermissions(calendar Calendar) error {
	return service.permissions.Add(
		authz.NewAllowPolicy(user.NewRef(calendar.ID), NewRef(calendar.ID), PermCalRead),
		authz.NewAllowPolicy(user.NewRef(calendar.ID), NewRef(calendar.ID), PermCalDelete),
	)
}
