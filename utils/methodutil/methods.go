package methodutil

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/vivasoft-ltd/go-ems/config"
	"github.com/vivasoft-ltd/go-ems/utils/errutil"
	"strconv"
)

func UserCacheKey(userID int) string {
	return config.Redis().MandatoryPrefix + config.Redis().UserPrefix + strconv.Itoa(userID)
}

func AccessUuidCacheKey(accessUuid string) string {
	return config.Redis().MandatoryPrefix + config.Redis().AccessUuidPrefix + accessUuid
}

func RefreshUuidCacheKey(refreshUuid string) string {
	return config.Redis().MandatoryPrefix + config.Redis().RefreshUuidPrefix + refreshUuid
}

func PermissionCacheKey(roleID int) string {
	return config.Redis().MandatoryPrefix + config.Redis().PermissionPrefix + strconv.Itoa(roleID)
}

func ParseJwtToken(token, secret string) (*jwt.Token, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errutil.ErrInvalidJwtSigningMethod
		}
		return []byte(secret), nil
	}

	return jwt.Parse(token, keyFunc)
}
