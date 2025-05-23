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
	"github.com/vivasoft-ltd/go-ems/utils/methodutil"
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
		logger.Error(fmt.Sprintf("error occurred: [%v] while generating password hash", err))
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
		logger.Error(fmt.Sprintf("error occurred: [%v] while creating user, email: [%s]", err, user.Email))
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
		logger.Error(fmt.Sprintf("error occurred: [%v] while updating user, user id: [%d]", err, user.ID))
		return err
	}

	go func() {
		if err := svc.redisSvc.Del(methodutil.UserCacheKey(existingUser.ID)); err != nil {
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
		if err := svc.redisSvc.Del(methodutil.UserCacheKey(existingUser.ID)); err != nil {
			logger.Error(err)
		}
	}()

	return nil
}

func (svc *UserServiceImpl) IsEmailExist(email string) (bool, error) {
	count, err := svc.repo.UserCountByEmail(email)
	if err != nil {
		logger.Error(fmt.Sprintf("error occurred: [%v] while fetching user by email", err))
		return false, err
	}

	return count != 0, nil

}

func (svc *UserServiceImpl) ReadUserByEmail(email string) (*models.User, error) {
	user, err := svc.repo.ReadUserByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error(fmt.Sprintf("error occurred: [%v] while fetching user by user email: [%s]", err, email))
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errutil.ErrUserNotFound
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
	err := svc.redisSvc.GetStruct(methodutil.UserCacheKey(id), &user)
	if err != nil && !errors.Is(err, redis.Nil) {
		logger.Error(fmt.Sprintf("error occurred: [%v] while fetching user from cache by user id: [%d]", err, id))
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
		logger.Error(fmt.Sprintf("error occurred: [%v] while fetching user by user id: [%d]", err, id))
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

func (svc *UserServiceImpl) StoreInCache(user *types.UserInfo) error {
	if err := svc.redisSvc.SetStruct(methodutil.UserCacheKey(user.ID), user, config.Redis().UserCacheTTL); err != nil {
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
		logger.Error(err)
		return nil, err
	}

	return permissions, nil
}

func (svc *UserServiceImpl) storePermissionInCache(roleID int, permissions []*models.Permission) error {
	if err := svc.redisSvc.SetStruct(methodutil.PermissionCacheKey(roleID), permissions, config.Redis().PermissionCacheTTL); err != nil {
		logger.Error(fmt.Sprintf("could not cache permissions in redis, role_id: [%d], err: [%v]", roleID, err))
		return err
	}
	return nil
}

func (svc *UserServiceImpl) ListUsers(req types.ListUserReq) (*types.PaginatedUserResp, error) {
	offset := (req.Page - 1) * req.Limit

	users, total, err := svc.repo.ReadPaginatedUsers(req.Limit, offset)
	if err != nil {
		logger.Error(fmt.Errorf("error occurred: [%v] while fetching paginated users", err))
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
