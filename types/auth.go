package types

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type (
	LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginResp struct {
		AccessToken  string    `json:"access_token"`
		RefreshToken string    `json:"refresh_token"`
		User         *UserInfo `json:"user"`
	}

	Token struct {
		UserID        int    `json:"uid"`
		AccessToken   string `json:"act"`
		RefreshToken  string `json:"rft"`
		AccessUuid    string `json:"aid"`
		RefreshUuid   string `json:"rid"`
		AccessExpiry  int64  `json:"axp"`
		RefreshExpiry int64  `json:"rxp"`
	}
)

func (l *LoginReq) Validate() error {
	return v.ValidateStruct(l,
		v.Field(&l.Email, v.Required, is.Email),
	)
}
