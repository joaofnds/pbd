package order

import (
	"app/authz"
	"app/internal/ref"
	"app/user"
)

const (
	PermOrderCreate = "order:create"
	PermOrderList   = "order:list"
	PermOrderRead   = "order:read"
	PermOrderUpdate = "order:update"
	PermOrderDelete = "order:delete"

	PermReviewCreate = "order:review:create"
	PermReviewRead   = "order:review:read"
	PermReviewDelete = "order:review:delete"
)

func NewRef(id string) ref.Ref {
	return ref.New("order", id)
}

type PermissionService struct {
	permissions authz.PermissionManager
}

func NewPermissionService(permissions authz.PermissionManager) *PermissionService {
	return &PermissionService{permissions: permissions}
}

func (service *PermissionService) GrantNewOrderPermissions(order Order) error {
	customer := user.NewRef(order.CustomerID)
	worker := user.NewRef(order.WorkerID)
	orderRef := NewRef(order.ID)

	return service.permissions.Add(
		authz.NewAllowPolicy(customer, orderRef, PermOrderRead),
		authz.NewAllowPolicy(customer, orderRef, PermOrderUpdate),
		authz.NewAllowPolicy(customer, orderRef, PermOrderDelete),
		authz.NewAllowPolicy(customer, orderRef, PermReviewRead),
		authz.NewAllowPolicy(customer, orderRef, PermReviewDelete),

		authz.NewAllowPolicy(worker, orderRef, PermOrderRead),
		authz.NewAllowPolicy(worker, orderRef, PermOrderUpdate),
		authz.NewAllowPolicy(worker, orderRef, PermOrderDelete),
		authz.NewAllowPolicy(worker, orderRef, PermReviewRead),
	)
}

func (service *PermissionService) GrantCompletedOrderPermissions(order Order) error {
	customer := user.NewRef(order.CustomerID)
	orderRef := NewRef(order.ID)

	return service.permissions.Add(
		authz.NewAllowPolicy(customer, orderRef, PermReviewCreate),
	)
}
