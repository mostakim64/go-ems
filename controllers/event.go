package controllers

import (
	"errors"
	"github.com/vivasoft-ltd/go-ems/consts"
	"github.com/vivasoft-ltd/go-ems/middlewares"
	"net/http"
	"strconv"

	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"github.com/vivasoft-ltd/go-ems/domain"
	"github.com/vivasoft-ltd/go-ems/types"
	"github.com/vivasoft-ltd/go-ems/utils/errutil"
	"github.com/vivasoft-ltd/go-ems/utils/msgutil"
	"github.com/vivasoft-ltd/golang-course-utils/logger"
)

type EventController struct {
	eventSvc domain.EventService
}

func NewEventController(eventSvc domain.EventService) *EventController {
	return &EventController{
		eventSvc: eventSvc,
	}
}

func (ctrl *EventController) CreateEvent(c echo.Context) error {
	var req types.CreateEventRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, msgutil.InvalidRequestMsg())
	}

	if err := req.Validate(); err != nil {
		logger.Error("validation error: %v", err)
		return c.JSON(http.StatusBadRequest, &types.ValidationError{
			Error: err,
		})
	}

	user, err := middlewares.CurrentUserFromCtx(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, msgutil.UserUnauthorized())
	}
	req.CreatedBy = user.ID

	resp, err := ctrl.eventSvc.CreateEvent(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, msgutil.SomethingWentWrongMsg())
	}

	return c.JSON(http.StatusOK, resp)
}

func (ctrl *EventController) ListEvents(c echo.Context) error {
	req := types.ListEventRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, msgutil.InvalidRequestMsg())
	}
	if req.Limit <= 0 {
		req.Limit = consts.DefaultPageSize
	}
	if req.Page <= 0 {
		req.Page = consts.DefaultPage
	}
	events, err := ctrl.eventSvc.ListEvents(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, msgutil.SomethingWentWrongMsg())
	}
	return c.JSON(http.StatusOK, events)
}

func (ctrl *EventController) ReadEventByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, msgutil.InvalidRequestMsg())
	}

	// Validate ID
	if err := v.Validate(id, v.Required); err != nil {
		logger.Error("validation error: %v", err)
		return c.JSON(http.StatusBadRequest, &types.ValidationError{
			Error: err,
		})
	}

	event, err := ctrl.eventSvc.ReadEventByID(id)
	if errors.Is(err, errutil.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, msgutil.EventNotFound())
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, msgutil.SomethingWentWrongMsg())
	}
	return c.JSON(http.StatusOK, event)
}

func (ctrl *EventController) UpdateEvent(c echo.Context) error {
	var req types.UpdateEventRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, msgutil.InvalidRequestMsg())
	}

	if err := req.Validate(); err != nil {
		logger.Error("validation error: %v", err)
		return c.JSON(http.StatusBadRequest, &types.ValidationError{
			Error: err,
		})
	}

	resp, err := ctrl.eventSvc.UpdateEvent(&req)
	if errors.Is(err, errutil.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, msgutil.EventNotFound())
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, msgutil.SomethingWentWrongMsg())
	}
	return c.JSON(http.StatusOK, resp)
}

func (ctrl *EventController) DeleteEvent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, msgutil.InvalidRequestMsg())
	}

	// Validate ID
	if err := v.Validate(id, v.Required); err != nil {
		logger.Error("validation error: %v", err)
		return c.JSON(http.StatusBadRequest, &types.ValidationError{
			Error: err,
		})
	}

	resp, err := ctrl.eventSvc.DeleteEvent(id)
	if errors.Is(err, errutil.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, msgutil.EventNotFound())
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, msgutil.SomethingWentWrongMsg())
	}
	return c.JSON(http.StatusNoContent, resp)
}
