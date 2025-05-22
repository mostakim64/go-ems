package domain

import "github.com/vivasoft-ltd/go-ems/types"

type (
	TokenService interface {
		CreateToken(userID int) (*types.Token, error)
		ParseAccessToken(accessToken string) (*types.Token, error)
		StoreTokenUUID(token *types.Token) error
		DeleteTokenUUID(token *types.Token) error
		ReadUserIDFromAccessTokenUUID(accessTokenUuid string) (int, error)
	}
)
