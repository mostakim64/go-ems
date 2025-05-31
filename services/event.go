package services

import (
	"errors"
	"fmt"
	"github.com/vivasoft-ltd/go-ems/consts"
	"github.com/vivasoft-ltd/go-ems/domain"
	"github.com/vivasoft-ltd/go-ems/models"
	"github.com/vivasoft-ltd/go-ems/types"
	"github.com/vivasoft-ltd/go-ems/utils/errutil"
)

type EventServiceImpl struct {
	eventRepo domain.EventRepository
	userRepo  domain.UserRepository
}

func NewEventServiceImpl(eventRepo domain.EventRepository, userRepo domain.UserRepository) *EventServiceImpl {
	return &EventServiceImpl{
		eventRepo: eventRepo,
		userRepo:  userRepo,
	}
}

func (svc *EventServiceImpl) CreateEvent(eventReq *types.CreateEventRequest) (*types.CreateEventResponse, error) {
	event := eventReq.ToEvent()
	if !eventReq.IsPublic && len(eventReq.Attendees) > 0 {
		users, err := svc.userRepo.ReadUsers(eventReq.Attendees)
		if err != nil {
			return nil, err
		}
		event.Attendees = users
	}

	createdEvent, err := svc.eventRepo.CreateEvent(event)
	if err != nil {
		return nil, err
	}

	return &types.CreateEventResponse{
		Message: "Event created",
		Event:   createdEvent,
	}, nil
}

func (svc *EventServiceImpl) ListEvents(user *types.CurrentUser) ([]*models.Event, error) {
	filter := svc.getEventListFilter(user)
	events, err := svc.eventRepo.ListEvents(filter)
	if errors.Is(err, errutil.ErrRecordNotFound) {
		return []*models.Event{}, nil
	}
	if err != nil {
		return nil, err
	}
	return events, nil
}
func (svc *EventServiceImpl) getEventListFilter(user *types.CurrentUser) *types.EventFilter {
	fmt.Printf("user: + %+v", user)
	filter := &types.EventFilter{}
	if user == nil {
		t := true
		filter.IsPublic = &t
		return filter
	}

	if user.HasPermission(consts.PermissionFetchAllEvent) {
		return filter
	}
	if user.HasPermission(consts.PermissionFetchOwnEvent) {
		filter.CreatedBy = &user.ID
	}
	if user.HasPermission(consts.PermissionFetchInvitedEvent) {
		filter.Attendee = &user.ID
	}

	return filter
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
func (svc *EventServiceImpl) RsvpEvent(request types.RsvpEventRequest) error {
	event, err := svc.eventRepo.ReadEventByID(request.EventID)
	if err != nil {
		return err
	}
	if !event.IsPublic {
		invitation, err := svc.eventRepo.ReadEventInvitation(event.ID, request.UserID)
		if invitation == nil || err != nil {
			return errutil.ErrRecordNotFound
		}
		invitation.StatusID = request.StatusID
		err = svc.eventRepo.UpsertEventInvitation(invitation)
		if err != nil {
			return err
		}
		return nil
	}
	count, err := svc.eventRepo.GetEventAttendeesCount(request.EventID)
	if err != nil {
		return err
	}
	if *event.Limit > 0 && count >= *event.Limit {
		return errutil.ErrEventCapacityExceeded
	}

	newInvitation := &models.EventAttendee{
		EventID:  request.EventID,
		UserID:   request.UserID,
		StatusID: request.StatusID,
	}
	err = svc.eventRepo.UpsertEventInvitation(newInvitation)
	if err != nil {
		return err
	}

	return nil
}
