package services

import (
	"github.com/vivasoft-ltd/go-ems/consts"
	"github.com/vivasoft-ltd/go-ems/domain"
	"github.com/vivasoft-ltd/go-ems/types"
	"github.com/vivasoft-ltd/go-ems/utils/errutil"
	"github.com/vivasoft-ltd/golang-course-utils/logger"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	userSvc  domain.UserService
	tokenSvc domain.TokenService
}

func NewAuthServiceImpl(userSvc domain.UserService, tokenSvc domain.TokenService) *AuthServiceImpl {
	return &AuthServiceImpl{userSvc: userSvc, tokenSvc: tokenSvc}
}

func (svc *AuthServiceImpl) Login(req *types.LoginReq) (*types.LoginResp, error) {
	user, err := svc.userSvc.ReadUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errutil.ErrInvalidLoginCredentials
	}

	token, err := svc.tokenSvc.CreateToken(user.ID)
	if err != nil {
		return nil, err
	}

	if err := svc.tokenSvc.StoreTokenUUID(token); err != nil {
		return nil, err
	}

	userInfo := &types.UserInfo{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		RoleID:    user.RoleID,
		Role:      consts.RoleMap[user.RoleID],
	}

	go func() {
		if err := svc.userSvc.StoreInCache(userInfo); err != nil {
			logger.Error(err)
		}
	}()

	resp := &types.LoginResp{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		User:         userInfo,
	}

	return resp, nil
}

func (svc *AuthServiceImpl) VerifyAccessToken(accessToken string) (*types.UserInfo, *types.Token, error) {
	token, err := svc.tokenSvc.ParseAccessToken(accessToken)
	if err != nil {
		return nil, nil, err
	}

	userID, err := svc.tokenSvc.ReadUserIDFromAccessTokenUUID(token.AccessUuid)
	if err != nil {
		return nil, nil, err
	}

	if userID != token.UserID {
		return nil, nil, errutil.ErrInvalidAccessToken
	}

	user, err := svc.userSvc.ReadUser(userID, true)
	if err != nil {
		return nil, nil, err
	}

	return user, token, nil
}

func (svc *AuthServiceImpl) Logout(accessTokenUuid, refreshTokenUuid string) error {
	return svc.tokenSvc.DeleteTokenUUID(&types.Token{AccessUuid: accessTokenUuid, RefreshUuid: refreshTokenUuid})
}
