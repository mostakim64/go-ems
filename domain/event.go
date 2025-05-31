package domain

import (
	"github.com/vivasoft-ltd/go-ems/models"
	"github.com/vivasoft-ltd/go-ems/types"
)

type (
	EventRepository interface {
		CreateEvent(event *models.Event) (*models.Event, error)
		ListEvents(filter *types.EventFilter) ([]*models.Event, error)
		ReadEventByID(id int) (*models.Event, error)
		UpdateEvent(event *models.Event) (*models.Event, error)
		DeleteEvent(id int) error
		ReadEventInvitation(eventID int, userID int) (*models.EventAttendee, error)
		UpsertEventInvitation(event *models.EventAttendee) error
		GetEventAttendeesCount(eventID int) (int, error)
	}

	EventService interface {
		CreateEvent(eventReq *types.CreateEventRequest) (*types.CreateEventResponse, error)
		ListEvents(user *types.CurrentUser) ([]*models.Event, error)
		ReadEventByID(id int) (*models.Event, error)
		DeleteEvent(id int) (*types.DeleteEventResponse, error)
		UpdateEvent(eventReq *types.UpdateEventRequest) (*types.UpdateEventResponse, error)
		RsvpEvent(request types.RsvpEventRequest) error
	}
)
