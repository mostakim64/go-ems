package types

import (
	"time"

	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/vivasoft-ltd/go-ems/models"
)

type (
	CreateEventRequest struct {
		Title       string  `json:"title"`
		Description *string `json:"description"`
		Location    *string `json:"location"`
		StartTime   *string `json:"start_time"`
		EndTime     *string `json:"end_time"`
		CreatedBy   int     `json:"created_by"`
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
)

func (cereq *CreateEventRequest) Validate() error {
	return v.ValidateStruct(cereq,
		v.Field(&cereq.Title, v.Required),
		v.Field(&cereq.Description, v.When(cereq.Description != nil, v.Length(0, 500))),
		v.Field(&cereq.Location, v.When(cereq.Location != nil, v.Length(0, 255))),
		v.Field(&cereq.StartTime, v.When(cereq.StartTime != nil, v.Date(time.RFC3339))),
		v.Field(&cereq.EndTime, v.When(cereq.EndTime != nil, v.Date(time.RFC3339))),
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
