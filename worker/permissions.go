package worker

import (
	"app/authz"
	"app/internal/ref"
	"app/user"
)

func NewRef(id string) ref.Ref {
	return ref.New("worker", id)
}

type PermissionService struct {
	permissions authz.PermissionManager
}

func NewPermissionServiceService(permissions authz.PermissionManager) *PermissionService {
	return &PermissionService{permissions: permissions}
}

func (service *PermissionService) GrantNewWorkerPermissions(worker Worker) error {
	return service.permissions.Add(
		authz.NewAllowPolicy(user.NewRef(worker.ID), NewRef(worker.ID), "*"),
	)
}
