package services

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/vivasoft-ltd/go-ems/config"
	"github.com/vivasoft-ltd/go-ems/consts"
	"github.com/vivasoft-ltd/go-ems/domain"
	"github.com/vivasoft-ltd/go-ems/models"
	"github.com/vivasoft-ltd/go-ems/types"
	"github.com/vivasoft-ltd/go-ems/utils/errutil"
	"github.com/vivasoft-ltd/golang-course-utils/logger"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strconv"
)

type UserServiceImpl struct {
	redisSvc *RedisService
	repo     domain.UserRepository
}

func NewUserServiceImpl(redisSvc *RedisService, userRepo domain.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		redisSvc: redisSvc,
		repo:     userRepo,
	}
}

func (svc *UserServiceImpl) CreateUser(req *types.CreateUserReq) error {
	isExist, err := svc.IsEmailExist(req.Email)
	if err != nil {
		return err
	}

	if isExist {
		return errutil.ErrUserAlreadyExist
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Email:     req.Email,
		Password:  string(hashedPass),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		RoleID:    req.RoleID,
	}

	if _, err := svc.repo.CreateUser(user); err != nil {
		return err
	}

	return nil
}

func (svc *UserServiceImpl) UpdateUser(req *types.UpdateUserReq) error {
	existingUser, err := svc.ReadUser(req.ID, false)
	if err != nil {
		return err
	}
	user := &models.User{
		ID:        existingUser.ID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := svc.repo.UpdateUser(user); err != nil {
		return err
	}

	go func() {
		if err := svc.redisSvc.Del(config.Redis().MandatoryPrefix + config.Redis().UserPrefix + strconv.Itoa(user.ID)); err != nil {
			logger.Error(err)
		}
	}()

	return nil
}

func (svc *UserServiceImpl) DeleteUser(id int) error {
	existingUser, err := svc.ReadUser(id, false)
	if err != nil {
		return err
	}

	if err := svc.repo.DeleteUser(existingUser.ID); err != nil {
		return err
	}

	go func() {
		if err := svc.redisSvc.Del(config.Redis().MandatoryPrefix + config.Redis().UserPrefix + strconv.Itoa(id)); err != nil {
			logger.Error(err)
		}
	}()

	return nil
}

func (svc *UserServiceImpl) IsEmailExist(email string) (bool, error) {
	count, err := svc.repo.UserCountByEmail(email)
	if err != nil {
		return false, err
	}

	return count != 0, nil

}

func (svc *UserServiceImpl) ReadUserByEmail(email string) (*models.User, error) {
	user, err := svc.repo.ReadUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (svc *UserServiceImpl) ReadUser(id int, fromCache bool) (*types.UserInfo, error) {
	if fromCache {
		return svc.readUserFromCache(id)
	}

	return svc.readUserFromDB(id)
}

func (svc *UserServiceImpl) readUserFromCache(id int) (*types.UserInfo, error) {
	var user *types.UserInfo
	err := svc.redisSvc.GetStruct(config.Redis().MandatoryPrefix+config.Redis().UserPrefix+strconv.Itoa(id), &user)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if errors.Is(err, redis.Nil) {
		return svc.readUserFromDB(id)
	}

	return user, nil
}

func (svc *UserServiceImpl) readUserFromDB(id int) (*types.UserInfo, error) {
	user, err := svc.repo.ReadUserById(id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errutil.ErrUserNotFound
	}

	return &types.UserInfo{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		RoleID:    user.RoleID,
		Role:      consts.RoleMap[user.RoleID],
	}, nil
}

// TODO: refactor redis key
func (svc *UserServiceImpl) StoreInCache(user *types.UserInfo) error {
	if err := svc.redisSvc.SetStruct(config.Redis().MandatoryPrefix+config.Redis().UserPrefix+strconv.Itoa(user.ID), user, config.Redis().UserCacheTTL); err != nil {
		logger.Error(fmt.Sprintf("could not cache user in redis, id: [%d], err: [%v]", user.ID, err))
	}
	return nil
}

func (svc *UserServiceImpl) ReadPermissionsByRole(roleID int) ([]*models.Permission, error) {
	permissions, err := svc.readPermissionFromCache(roleID)
	if err != nil {
		return nil, err
	}

	if permissions != nil {
		return permissions, nil
	}

	permissions, err = svc.repo.ReadPermissionsByRole(roleID)
	if err != nil {
		return nil, err
	}

	if err := svc.storePermissionInCache(roleID, permissions); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (svc *UserServiceImpl) readPermissionFromCache(roleID int) ([]*models.Permission, error) {
	var permissions []*models.Permission
	err := svc.redisSvc.GetStruct(config.Redis().MandatoryPrefix+config.Redis().PermissionPrefix+strconv.Itoa(roleID), &permissions)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	return permissions, nil
}

func (svc *UserServiceImpl) storePermissionInCache(roleID int, permissions []*models.Permission) error {
	if err := svc.redisSvc.SetStruct(config.Redis().MandatoryPrefix+config.Redis().PermissionPrefix+strconv.Itoa(roleID), permissions, config.Redis().PermissionCacheTTL); err != nil {
		logger.Error(fmt.Sprintf("could not cache permissions in redis, role_id: [%d], err: [%v]", roleID, err))
		return err
	}
	return nil
}

//	func (service *Service) UpdateUser(req types.UserUpdateReq) error {
//		userUpdReq := domain.User{
//			ID:        req.ID,
//			UserName:  req.UserName,
//			FirstName: req.FirstName,
//			LastName:  req.LastName,
//			Phone:     req.Phone,
//		}
//
//		err := service.repo.UpdateUser(userUpdReq)
//		if err != nil {
//			slog.Error("failed to update user", "error", err)
//			return fmt.Errorf("failed to update user: %w", err)
//		}
//		return nil
//	}
//
//	func (service *Service) ReadUserByEmail(email string) (types.UserRolePermissionsInfo, error) {
//		userPermissions, err := service.repo.GetUserPermissionsByEmail(email)
//		if err != nil {
//			return types.UserRolePermissionsInfo{}, err
//		}
//		if len(userPermissions) <= 0 {
//			return types.UserRolePermissionsInfo{}, nil
//		}
//
//		user := userPermissions[0]
//
//		resp := types.UserRolePermissionsInfo{
//			ID:       user.ID,
//			UserName: user.UserName,
//			Email:    user.Email,
//			Role:     user.RoleName,
//			Password: user.Password,
//		}
//
//		permissions := []string{}
//		for _, userPerm := range userPermissions {
//			permissions = append(permissions, userPerm.PermissionName)
//		}
//
//		resp.Permissions = permissions
//
//		return resp, nil
//	}
//
//	func (service *Service) AssignRoleToUser(req types.AssignRoleRequest) error {
//		userRoleReq := []domain.UserRole{}
//		for _, roleID := range req.RoleIDs {
//			userRoleReq = append(userRoleReq, domain.UserRole{
//				UserID: req.ID,
//				RoleID: roleID,
//			})
//		}
//
//		return service.repo.AssignRoleToUser(userRoleReq)
//	}
//
//	func (service *Service) GetUserRoles(userID int) (types.UserRolesResp, error) {
//		response := types.UserRolesResp{}
//		userRoles, err := service.repo.GetUserRoles([]int{userID})
//		if err != nil {
//			return response, err
//		}
//
//		if len(userRoles) <= 0 {
//			return response, nil
//		}
//		response.ID = userRoles[0].UserID
//		response.Roles = make([]types.RoleResponse, len(userRoles))
//
//		for i, userRole := range userRoles {
//			response.Roles[i] = types.RoleResponse{
//				ID:          userRole.UserID,
//				Name:        userRole.RoleName,
//				Description: userRole.RoleDescription,
//			}
//		}
//
//		return response, nil
//	}
//
// func (service *Service) DeleteUserRole(req types.DeleteUserRoleReq) error {
//
//		deleteUserRoleReq := domain.DeleteUserRoleReq{
//			UserID: req.ID,
//			RoleID: req.RoleID,
//		}
//		err := service.repo.DeleteUserRole(deleteUserRoleReq)
//		if err != nil {
//			slog.Error("failed to delete user role", "error", err)
//			return fmt.Errorf("failed to delete user role: %w", err)
//		}
//		return nil
//	}
//
//	func (service *Service) DeleteUser(userID int) error {
//		err := service.repo.DeleteUser(userID)
//		if err != nil {
//			slog.Error("failed to delete user", "error", err)
//			return fmt.Errorf("failed to delete user: %w", err)
//		}
//		return nil
//	}
//
//	func (service *Service) ReadUser(userID int) (types.UserResp, error) {
//		user, err := service.repo.ReadUser(userID)
//		if err != nil {
//			slog.Error("failed to get user", "error", err)
//			return types.UserResp{}, fmt.Errorf("failed to get user: %w", err)
//		}
//		resp := mapToUserResponse(user)
//		return resp, nil
//	}
func (svc *UserServiceImpl) ListUsers(req types.ListUserReq) (*types.PaginatedUserResp, error) {
	offset := (req.Page - 1) * req.Limit

	users, total, err := svc.repo.ReadPaginatedUsers(req.Limit, offset)
	if err != nil {
		return nil, err
	}

	for i, user := range users {
		users[i].Role = consts.RoleMap[user.RoleID]
	}

	resp := &types.PaginatedUserResp{
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
		Users: users,
	}

	return resp, nil
}

//func (service *Service) GetUserPermissionsByID(userID int) (types.UserWithParamsResp, error) {
//	userPermissions, err := service.repo.GetUserPermissionsByID(userID)
//	if err != nil {
//		return types.UserWithParamsResp{}, err
//	}
//	if len(userPermissions) <= 0 {
//		return types.UserWithParamsResp{}, nil
//	}
//
//	user := userPermissions[0]
//
//	resp := types.UserWithParamsResp{
//		CurrentUser: types.CurrentUser{
//			ID:       user.ID,
//			Email:    user.Email,
//			UserName: user.UserName,
//		},
//		Role: user.RoleName,
//	}
//
//	permissions := []string{}
//	for _, userPerm := range userPermissions {
//		permissions = append(permissions, userPerm.PermissionName)
//	}
//
//	resp.Permissions = permissions
//
//	return resp, nil
//}
//
//func (service *Service) ResetPassword(req types.ResetPasswordReq) error {
//
//	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
//	if err != nil {
//		return err
//	}
//
//	resetPasswordReq := domain.ResetPasswordReq{
//		ID:       req.ID,
//		Password: string(hashedPass),
//	}
//	err = service.repo.ResetPassword(resetPasswordReq)
//	if err != nil {
//		slog.Error("failed to reset password", "error", err)
//		return fmt.Errorf("failed to reset password: %w", err)
//	}
//	return nil
//}
//
//func mapToUserResponse(user domain.UserResp) types.UserResp {
//	resp := types.UserResp{
//		ID:        user.ID,
//		UserName:  user.UserName,
//		FirstName: user.FirstName,
//		LastName:  user.LastName,
//		Email:     user.Email,
//		Phone:     user.Phone,
//		Password:  user.Password,
//	}
//	return resp
//}
