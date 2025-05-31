package types

import (
	"time"

	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/vivasoft-ltd/go-ems/models"
)

type (
	CreateEventRequest struct {
		Title         string  `json:"title"`
		Description   *string `json:"description"`
		Location      *string `json:"location"`
		StartTime     *string `json:"start_time"`
		EndTime       *string `json:"end_time"`
		CreatedBy     int     `json:"created_by"`
		IsPublic      bool    `json:"is_public"`
		AttendeeLimit *int    `json:"attendee_limit"`
		Attendees     []int   `json:"attendees"`
	}

	UpdateEventRequest struct {
		ID int `param:"id"`
		CreateEventRequest
	}

	CreateEventResponse struct {
		Message string        `json:"message"`
		Event   *models.Event `json:"event"`
	}

	DeleteEventResponse struct {
		Message string `json:"message"`
	}

	UpdateEventResponse struct {
		Message string        `json:"message"`
		Event   *models.Event `json:"event"`
	}
	RsvpEventRequest struct {
		EventID  int `json:"event_id" param:"id"`
		UserID   int `json:"user_id"`
		StatusID int `json:"status_id"`
	}
	EventFilter struct {
		CreatedBy *int  `query:"created_by"`
		Attendee  *int  `query:"attendee"`
		IsPublic  *bool `query:"is_public"`
	}
	ListEventRequest struct {
		Page  int `query:"page"`
		Limit int `query:"limit"`
	}
	PaginatedEventResponse struct {
		Total  int             `json:"total"`
		Page   int             `json:"page"`
		Limit  int             `json:"limit"`
		Events []*models.Event `json:"events"`
	}
)

func (r *RsvpEventRequest) Validate() error {
	return v.ValidateStruct(r,
		v.Field(&r.EventID, v.Required),
		v.Field(&r.StatusID, v.Required, v.In(2, 3)),
	)
}

func (cereq *CreateEventRequest) Validate() error {
	return v.ValidateStruct(cereq,
		v.Field(&cereq.Title, v.Required),
		v.Field(&cereq.Description, v.When(cereq.Description != nil, v.Length(0, 500))),
		v.Field(&cereq.Location, v.When(cereq.Location != nil, v.Length(0, 255))),
		v.Field(&cereq.StartTime, v.When(cereq.StartTime != nil, v.Date(time.RFC3339))),
		v.Field(&cereq.EndTime, v.When(cereq.EndTime != nil, v.Date(time.RFC3339))),
		v.Field(&cereq.Attendees, v.When(!cereq.IsPublic, v.Required, v.Length(1, 0))),
	)
}

func (uereq *UpdateEventRequest) Validate() error {
	return v.ValidateStruct(uereq,
		v.Field(&uereq.ID, v.Required),
		v.Field(&uereq.CreateEventRequest, v.Required),
	)
}

func (cereq *CreateEventRequest) ToEvent() *models.Event {
	event := &models.Event{
		Title:       cereq.Title,
		Description: cereq.Description,
		Location:    cereq.Location,
		CreatedBy:   cereq.CreatedBy,
		IsPublic:    cereq.IsPublic,
		Limit:       cereq.AttendeeLimit,
	}
	if cereq.StartTime != nil {
		event.StartTime, _ = parseTime(*cereq.StartTime, time.RFC3339)
	}
	if cereq.EndTime != nil {
		event.EndTime, _ = parseTime(*cereq.EndTime, time.RFC3339)
	}
	return event
}

func (uereq *UpdateEventRequest) ToEvent() *models.Event {
	event := &models.Event{
		ID:          uereq.ID,
		Title:       uereq.Title,
		Description: uereq.Description,
		Location:    uereq.Location,
		CreatedBy:   uereq.CreatedBy,
	}
	if uereq.StartTime != nil {
		event.StartTime, _ = parseTime(*uereq.StartTime, time.RFC3339)
	}
	if uereq.EndTime != nil {
		event.EndTime, _ = parseTime(*uereq.EndTime, time.RFC3339)
	}
	return event
}
