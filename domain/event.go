package domain

import (
	"github.com/vivasoft-ltd/go-ems/models"
	"github.com/vivasoft-ltd/go-ems/types"
)

type (
	EventRepository interface {
		CreateEvent(event *models.Event) (*models.Event, error)
		ListEvents(limit, offset int) ([]*models.Event, int, error)
		ReadEventByID(id int) (*models.Event, error)
		UpdateEvent(event *models.Event) (*models.Event, error)
		DeleteEvent(id int) error
	}

	EventService interface {
		CreateEvent(eventReq *types.CreateEventRequest) (*types.CreateEventResponse, error)
		ListEvents(req types.ListEventRequest) (*types.PaginatedEventResponse, error)
		ReadEventByID(id int) (*models.Event, error)
		DeleteEvent(id int) (*types.DeleteEventResponse, error)
		UpdateEvent(eventReq *types.UpdateEventRequest) (*types.UpdateEventResponse, error)
	}
)
