package domain

import (
	"github.com/vivasoft-ltd/go-ems/models"
	"github.com/vivasoft-ltd/go-ems/types"
)

type (
	UserService interface {
		CreateUser(req *types.CreateUserReq) error
		UpdateUser(req *types.UpdateUserReq) error
		ReadUser(id int, fromCache bool) (*types.UserInfo, error)
		DeleteUser(id int) error
		ReadUserByEmail(email string) (*models.User, error)
		StoreInCache(user *types.UserInfo) error
		ListUsers(req types.ListUserReq) (*types.PaginatedUserResp, error)
		ReadPermissionsByRole(roleID int) ([]*models.Permission, error)
		ListAttendees(user types.CurrentUser) ([]types.AttendeeResp, error)
	}
	UserRepository interface {
		CreateUser(user *models.User) (*models.User, error)
		ReadUserById(id int) (*models.User, error)
		ReadUsers(id []int) ([]models.User, error)
		ReadPaginatedUsers(limit, offset int) ([]*types.UserInfo, int, error)
		UpdateUser(user *models.User) error
		DeleteUser(id int) error
		ReadUserByEmail(email string) (*models.User, error)
		UserCountByEmail(email string) (int, error)
		ReadPermissionsByRole(roleID int) ([]*models.Permission, error)
		ListAttendees(filter *types.AttendeeFilter) ([]types.AttendeeResp, error)
	}
)
