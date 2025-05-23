package types

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/vivasoft-ltd/go-ems/consts"
)

type (
	CurrentUser struct {
		ID          int    `json:"id"`
		Email       string `json:"email"`
		RoleID      int    `json:"role_id"`
		Role        string `json:"role"`
		AccessUuid  string `json:"access_uuid"`
		RefreshUuid string `json:"refresh_uuid"`
	}

	CreateUserReq struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		RoleID    int    `json:"role_id"`
	}

	UpdateUserReq struct {
		ID        int    `param:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		RoleID    int    `json:"role_id"`
	}

	UserReq struct {
		ID int `param:"id"`
	}

	UserInfo struct {
		ID        int    `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		RoleID    int    `json:"role_id"`
		Role      string `json:"role,omitempty" gorm:"-"`
	}

	ListUserReq struct {
		Page  int `query:"page"`
		Limit int `query:"limit"`
	}

	PaginatedUserResp struct {
		Total int         `json:"total"`
		Page  int         `json:"page"`
		Limit int         `json:"limit"`
		Users []*UserInfo `json:"users"`
	}
)

func (crq *CreateUserReq) Validate() error {
	return v.ValidateStruct(crq,
		v.Field(&crq.Email, v.Required, is.Email),
		v.Field(&crq.Password, v.Required),
		v.Field(&crq.FirstName, v.Required, v.Length(0, 50)),
		v.Field(&crq.LastName, v.Required, v.Length(0, 50)),
		v.Field(&crq.RoleID, v.Required, v.In(consts.RoleIdAdmin, consts.RoleIdManager, consts.RoleIdAttendee)),
	)
}

func (urq *UpdateUserReq) Validate() error {
	return v.ValidateStruct(urq,
		v.Field(&urq.FirstName, v.Required, v.Length(0, 50)),
		v.Field(&urq.LastName, v.Required, v.Length(0, 50)),
		v.Field(&urq.RoleID, v.Required, v.In(consts.RoleIdAdmin, consts.RoleIdManager, consts.RoleAttendee)),
	)
}

func (rq *UserReq) Validate() error {
	return v.ValidateStruct(rq,
		v.Field(&rq.ID, v.Required, v.Min(1)),
	)
}
