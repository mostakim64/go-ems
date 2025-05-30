package services

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/vivasoft-ltd/go-ems/config"
	"github.com/vivasoft-ltd/go-ems/types"
	"github.com/vivasoft-ltd/go-ems/utils/errutil"
	"github.com/vivasoft-ltd/go-ems/utils/methodutil"
	"github.com/vivasoft-ltd/golang-course-utils/logger"
	"time"
)

type TokenServiceImpl struct {
	redisSvc *RedisService
}

func NewTokenServiceImpl(redisSvc *RedisService) *TokenServiceImpl {
	return &TokenServiceImpl{redisSvc: redisSvc}
}

func (*TokenServiceImpl) CreateToken(userID int) (*types.Token, error) {
	jwtConf := config.Jwt()
	token := &types.Token{}

	token.UserID = userID
	token.AccessExpiry = time.Now().Add(time.Second * jwtConf.AccessTokenExpiry).Unix()
	token.AccessUuid = uuid.New().String()

	token.RefreshExpiry = time.Now().Add(time.Second * jwtConf.RefreshTokenExpiry).Unix()
	token.RefreshUuid = uuid.New().String()

	atClaims := jwt.MapClaims{}
	atClaims["uid"] = userID
	atClaims["aid"] = token.AccessUuid
	atClaims["rid"] = token.RefreshUuid
	atClaims["exp"] = token.AccessExpiry

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	var err error
	token.AccessToken, err = at.SignedString([]byte(jwtConf.AccessTokenSecret))
	if err != nil {
		logger.Error(err)
		return nil, errutil.ErrAccessTokenSign
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["uid"] = userID
	rtClaims["aid"] = token.AccessUuid
	rtClaims["rid"] = token.RefreshUuid
	rtClaims["exp"] = token.RefreshExpiry

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	token.RefreshToken, err = rt.SignedString([]byte(jwtConf.RefreshTokenSecret))
	if err != nil {
		log.Error(err)
		return nil, errutil.ErrRefreshTokenSign
	}

	return token, nil
}

func (svc *TokenServiceImpl) ParseAccessToken(accessToken string) (*types.Token, error) {
	parsedToken, err := methodutil.ParseJwtToken(accessToken, config.Jwt().AccessTokenSecret)
	if err != nil {
		log.Error(err)
		return nil, errutil.ErrParseJwt
	}

	if _, ok := parsedToken.Claims.(jwt.Claims); !ok || !parsedToken.Valid {
		return nil, errutil.ErrInvalidAccessToken
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errutil.ErrInvalidAccessToken
	}

	return mapClaimsToToken(claims)
}

func (svc *TokenServiceImpl) StoreTokenUUID(token *types.Token) error {
	err := svc.redisSvc.Set(methodutil.AccessUuidCacheKey(token.AccessUuid), token.UserID, time.Duration(token.AccessExpiry))
	if err != nil {
		return err
	}

	err = svc.redisSvc.Set(methodutil.RefreshUuidCacheKey(token.RefreshUuid), token.UserID, time.Duration(token.RefreshExpiry))
	if err != nil {
		return err
	}

	return nil
}

func (svc *TokenServiceImpl) DeleteTokenUUID(token *types.Token) error {
	err := svc.redisSvc.Del(methodutil.AccessUuidCacheKey(token.AccessUuid), methodutil.RefreshUuidCacheKey(token.RefreshUuid))

	if err != nil {
		return err
	}

	return nil
}

func (svc *TokenServiceImpl) ReadUserIDFromAccessTokenUUID(accessTokenUuid string) (int, error) {
	userID, err := svc.redisSvc.GetInt(methodutil.AccessUuidCacheKey(accessTokenUuid))
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func mapClaimsToToken(claims jwt.MapClaims) (*types.Token, error) {
	jsonData, err := json.Marshal(claims)
	if err != nil {
		return nil, err
	}

	var token types.Token
	err = json.Unmarshal(jsonData, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
