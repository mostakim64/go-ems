package services

import (
	"errors"

	"github.com/vivasoft-ltd/go-ems/domain"
	"github.com/vivasoft-ltd/go-ems/models"
	"github.com/vivasoft-ltd/go-ems/types"
	"github.com/vivasoft-ltd/go-ems/utils/errutil"
)

type EventServiceImpl struct {
	eventRepo domain.EventRepository
}

func NewEventServiceImpl(eventRepo domain.EventRepository) *EventServiceImpl {
	return &EventServiceImpl{
		eventRepo: eventRepo,
	}
}

func (svc *EventServiceImpl) CreateEvent(eventReq *types.CreateEventRequest) (*types.CreateEventResponse, error) {
	event := eventReq.ToEvent()
	createdEvent, err := svc.eventRepo.CreateEvent(event)
	if err != nil {
		return nil, err
	}

	return &types.CreateEventResponse{
		Message: "Event created",
		Event:   createdEvent,
	}, nil
}

func (svc *EventServiceImpl) ListEvents(req types.ListEventRequest) (*types.PaginatedEventResponse, error) {
	offset := (req.Page - 1) * req.Limit
	events, count, err := svc.eventRepo.ListEvents(req.Limit, offset)
	if errors.Is(err, errutil.ErrRecordNotFound) {
		return &types.PaginatedEventResponse{}, nil
	}
	if err != nil {
		return nil, err
	}
	response := &types.PaginatedEventResponse{
		Page:   req.Page,
		Limit:  req.Limit,
		Total:  count,
		Events: events,
	}
	return response, nil
}

func (svc *EventServiceImpl) ReadEventByID(id int) (*models.Event, error) {
	event, err := svc.eventRepo.ReadEventByID(id)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (svc *EventServiceImpl) UpdateEvent(eventReq *types.UpdateEventRequest) (*types.UpdateEventResponse, error) {
	existingEvent, err := svc.eventRepo.ReadEventByID(eventReq.ID)
	if err != nil {
		return nil, err
	}
	if existingEvent == nil {
		return nil, errutil.ErrRecordNotFound
	}

	event := eventReq.ToEvent()
	updatedEvent, err := svc.eventRepo.UpdateEvent(event)
	if err != nil {
		return nil, err
	}
	return &types.UpdateEventResponse{
		Message: "Event updated",
		Event:   updatedEvent,
	}, nil
}

func (svc *EventServiceImpl) DeleteEvent(id int) (*types.DeleteEventResponse, error) {
	err := svc.eventRepo.DeleteEvent(id)
	if err != nil {
		return nil, err
	}
	return &types.DeleteEventResponse{
		Message: "Event deleted",
	}, nil
}
