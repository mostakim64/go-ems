package controllers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/vivasoft-ltd/go-ems/domain"
	"github.com/vivasoft-ltd/go-ems/middlewares"
	"github.com/vivasoft-ltd/go-ems/types"
	"github.com/vivasoft-ltd/go-ems/utils/errutil"
	"github.com/vivasoft-ltd/go-ems/utils/msgutil"
	"net/http"
)

type AuthController struct {
	authSvc domain.AuthService
}

func NewAuthController(authSvc domain.AuthService) *AuthController {
	return &AuthController{authSvc: authSvc}
}

func (ctrl *AuthController) Login(c echo.Context) error {
	var req types.LoginReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, msgutil.InvalidRequestMsg())
	}

	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, &types.ValidationError{
			Error: err,
		})
	}

	resp, err := ctrl.authSvc.Login(&req)
	if err != nil {
		switch {
		case errors.Is(err, errutil.ErrUserNotFound):
		case errors.Is(err, errutil.ErrInvalidLoginCredentials):
			return c.JSON(http.StatusUnauthorized, msgutil.InvalidLoginCredentials())
		}
		return c.JSON(http.StatusInternalServerError, msgutil.SomethingWentWrongMsg())
	}

	return c.JSON(http.StatusOK, resp)
}

func (ctrl *AuthController) Logout(c echo.Context) error {
	user, err := middlewares.CurrentUserFromCtx(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, msgutil.UserUnauthorized())
	}

	if err := ctrl.authSvc.Logout(user.AccessUuid, user.RefreshUuid); err != nil {
		return c.JSON(http.StatusInternalServerError, msgutil.SomethingWentWrongMsg())
	}

	return c.NoContent(http.StatusOK)
}
