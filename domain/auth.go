package domain

import (
	"github.com/vivasoft-ltd/go-ems/types"
)

type (
	AuthService interface {
		Login(req *types.LoginReq) (*types.LoginResp, error)
		VerifyAccessToken(accessToken string) (*types.UserInfo, *types.Token, error)
		Logout(accessTokenUuid, refreshTokenUuid string) error
	}
)
