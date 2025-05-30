package controllers

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/vivasoft-ltd/go-ems/consts"
	"github.com/vivasoft-ltd/go-ems/domain"
	"github.com/vivasoft-ltd/go-ems/middlewares"
	"github.com/vivasoft-ltd/go-ems/types"
	"github.com/vivasoft-ltd/go-ems/utils/errutil"
	"github.com/vivasoft-ltd/go-ems/utils/msgutil"
	"net/http"
)

type UserController struct {
	userSvc domain.UserService
}

func NewUserController(userSvc domain.UserService) *UserController {
	return &UserController{userSvc: userSvc}
}

func (ctrl *UserController) Signup(c echo.Context) error {
	var req types.CreateUserReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, msgutil.InvalidRequestMsg())
	}

	req.RoleID = consts.RoleIdAttendee
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, &types.ValidationError{
			Error: err,
		})
	}

	if err := ctrl.userSvc.CreateUser(&req); err != nil {
		switch {
		case errors.Is(err, errutil.ErrUserAlreadyExist):
			return c.JSON(http.StatusConflict, msgutil.UserAlreadyExists())
		}
		return c.JSON(http.StatusInternalServerError, msgutil.SomethingWentWrongMsg())
	}

	return c.JSON(http.StatusCreated, msgutil.UserCreatedSuccessfully())
}

func (ctrl *UserController) Profile(c echo.Context) error {
	user, err := middlewares.CurrentUserFromCtx(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, msgutil.UserUnauthorized())
	}

	resp, err := ctrl.userSvc.ReadUser(user.ID, false)
	if err != nil {
		switch {
		case errors.Is(err, errutil.ErrUserNotFound):
			return c.JSON(http.StatusNotFound, msgutil.UserNotFound())
		}
		return c.JSON(http.StatusInternalServerError, msgutil.SomethingWentWrongMsg())
	}

	return c.JSON(http.StatusOK, resp)
}

func (ctrl *UserController) CreateUser(c echo.Context) error {
	fmt.Println("here")
	var req types.CreateUserReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, msgutil.InvalidRequestMsg())
	}

	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, &types.ValidationError{
			Error: err,
		})
	}

	if err := ctrl.userSvc.CreateUser(&req); err != nil {
		switch {
		case errors.Is(err, errutil.ErrUserAlreadyExist):
			return c.JSON(http.StatusConflict, msgutil.UserAlreadyExists())
		}
		return c.JSON(http.StatusInternalServerError, msgutil.SomethingWentWrongMsg())
	}
	return c.JSON(http.StatusCreated, msgutil.UserCreatedSuccessfully())
}

func (ctrl *UserController) UpdateUser(c echo.Context) error {
	var req types.UpdateUserReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, msgutil.InvalidRequestMsg())
	}

	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, &types.ValidationError{
			Error: err,
		})
	}
	if err := ctrl.userSvc.UpdateUser(&req); err != nil {
		switch {
		case errors.Is(err, errutil.ErrUserNotFound):
			return c.JSON(http.StatusNotFound, msgutil.UserNotFound())
		}
		return c.JSON(http.StatusInternalServerError, msgutil.SomethingWentWrongMsg())
	}

	return c.JSON(http.StatusOK, msgutil.UserUpdatedSuccessfully())
}

func (ctrl *UserController) DeleteUser(c echo.Context) error {
	var req types.UserReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, msgutil.InvalidRequestMsg())
	}

	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, &types.ValidationError{
			Error: err,
		})
	}

	if err := ctrl.userSvc.DeleteUser(req.ID); err != nil {
		switch {
		case errors.Is(err, errutil.ErrUserNotFound):
			return c.JSON(http.StatusNotFound, msgutil.UserNotFound())
		}
		return c.JSON(http.StatusInternalServerError, msgutil.SomethingWentWrongMsg())
	}

	return c.JSON(http.StatusOK, msgutil.UserDeletedSuccessfully())
}

func (ctrl *UserController) ReadUser(c echo.Context) error {
	var req types.UserReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, msgutil.InvalidRequestMsg())
	}

	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, &types.ValidationError{
			Error: err,
		})
	}

	user, err := ctrl.userSvc.ReadUser(req.ID, false)
	if err != nil {
		switch {
		case errors.Is(err, errutil.ErrUserNotFound):
			return c.JSON(http.StatusNotFound, msgutil.UserNotFound())
		}
		return c.JSON(http.StatusInternalServerError, msgutil.SomethingWentWrongMsg())
	}

	return c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) ListUsers(c echo.Context) error {
	var req types.ListUserReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, msgutil.InvalidRequestMsg())
	}

	if req.Limit <= 0 {
		req.Limit = consts.DefaultPageSize
	}

	if req.Page <= 0 {
		req.Page = consts.DefaultPage
	}
	resp, err := ctrl.userSvc.ListUsers(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, msgutil.SomethingWentWrongMsg())
	}

	return c.JSON(http.StatusOK, resp)
}

func (ctrl *UserController) ListAttendees(c echo.Context) error {
	user, err := middlewares.CurrentUserFromCtx(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, msgutil.UserUnauthorized())
	}
	users, err := ctrl.userSvc.ListAttendees(*user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, msgutil.SomethingWentWrongMsg())
	}

	return c.JSON(http.StatusOK, users)
}
