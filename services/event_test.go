package services

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/vivasoft-ltd/go-ems/consts"
	"github.com/vivasoft-ltd/go-ems/models"
	"github.com/vivasoft-ltd/go-ems/services/mocks"
	"github.com/vivasoft-ltd/go-ems/types"
	"github.com/vivasoft-ltd/go-ems/utils/errutil"
	"go.uber.org/mock/gomock"
)

// Helper functions for creating test data
func createTestEvent(id int) *models.Event {
	description := "Test Description"
	location := "Test Location"
	startTime := time.Now().Add(24 * time.Hour)
	endTime := time.Now().Add(48 * time.Hour)
	limit := 10

	return &models.Event{
		ID:          id,
		Title:       "Test Event",
		Description: &description,
		Location:    &location,
		StartTime:   &startTime,
		EndTime:     &endTime,
		IsPublic:    true,
		Limit:       &limit,
		CreatedBy:   1,
	}
}

// Test cases for EventServiceImpl.CreateEvent
func TestCreateEvent(t *testing.T) {
	// Test case 1: Successful creation of a public event
	t.Run("SuccessfulPublicEventCreation", func(t *testing.T) {
		t.Parallel() // Enable parallel execution for this subtest

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEventRepo := mocks.NewMockEventRepository(ctrl)
		mockUserRepo := mocks.NewMockUserRepository(ctrl)

		description := "Test Description"
		location := "Test Location"
		startTime := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
		endTime := time.Now().Add(48 * time.Hour).Format(time.RFC3339)
		limit := 10

		request := &types.CreateEventRequest{
			Title:         "Test Event",
			Description:   &description,
			Location:      &location,
			StartTime:     &startTime,
			EndTime:       &endTime,
			CreatedBy:     1,
			IsPublic:      true,
			AttendeeLimit: &limit,
		}

		expectedEvent := request.ToEvent()
		expectedEvent.ID = 1

		mockEventRepo.EXPECT().
			CreateEvent(gomock.Any()).
			Return(expectedEvent, nil)

		service := NewEventServiceImpl(mockEventRepo, mockUserRepo)
		response, err := service.CreateEvent(request)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if response == nil {
			t.Error("Expected response, got nil")
		}
		if response.Message != "Event created" {
			t.Errorf("Expected message 'Event created', got '%s'", response.Message)
		}
		if !reflect.DeepEqual(response.Event, expectedEvent) {
			t.Errorf("Expected event %v, got %v", expectedEvent, response.Event)
		}
	})

	// Test case 2: Successful creation of a private event with attendees
	t.Run("SuccessfulPrivateEventCreation", func(t *testing.T) {
		t.Parallel() // Enable parallel execution for this subtest

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEventRepo := mocks.NewMockEventRepository(ctrl)
		mockUserRepo := mocks.NewMockUserRepository(ctrl)

		description := "Test Description"
		location := "Test Location"
		startTime := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
		endTime := time.Now().Add(48 * time.Hour).Format(time.RFC3339)
		limit := 10

		request := &types.CreateEventRequest{
			Title:         "Test Event",
			Description:   &description,
			Location:      &location,
			StartTime:     &startTime,
			EndTime:       &endTime,
			CreatedBy:     1,
			IsPublic:      false,
			AttendeeLimit: &limit,
			Attendees:     []int{2, 3},
		}

		users := []models.User{
			{ID: 2, Email: "user2@example.com", FirstName: "User", LastName: "Two"},
			{ID: 3, Email: "user3@example.com", FirstName: "User", LastName: "Three"},
		}

		expectedEvent := request.ToEvent()
		expectedEvent.ID = 1
		expectedEvent.Attendees = users

		mockUserRepo.EXPECT().
			ReadUsers(gomock.Eq([]int{2, 3})).
			Return(users, nil)

		mockEventRepo.EXPECT().
			CreateEvent(gomock.Any()).
			Return(expectedEvent, nil)

		service := NewEventServiceImpl(mockEventRepo, mockUserRepo)
		response, err := service.CreateEvent(request)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if response == nil {
			t.Error("Expected response, got nil")
		}
		if response.Message != "Event created" {
			t.Errorf("Expected message 'Event created', got '%s'", response.Message)
		}
		if !reflect.DeepEqual(response.Event, expectedEvent) {
			t.Errorf("Expected event %v, got %v", expectedEvent, response.Event)
		}
	})

	// Test case 3: Error when reading users
	t.Run("ErrorReadingUsers", func(t *testing.T) {
		t.Parallel() // Enable parallel execution for this subtest

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEventRepo := mocks.NewMockEventRepository(ctrl)
		mockUserRepo := mocks.NewMockUserRepository(ctrl)

		description := "Test Description"
		location := "Test Location"
		startTime := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
		endTime := time.Now().Add(48 * time.Hour).Format(time.RFC3339)
		limit := 10

		request := &types.CreateEventRequest{
			Title:         "Test Event",
			Description:   &description,
			Location:      &location,
			StartTime:     &startTime,
			EndTime:       &endTime,
			CreatedBy:     1,
			IsPublic:      false,
			AttendeeLimit: &limit,
			Attendees:     []int{2, 3},
		}

		mockUserRepo.EXPECT().
			ReadUsers(gomock.Eq([]int{2, 3})).
			Return(nil, errors.New("error reading users"))

		service := NewEventServiceImpl(mockEventRepo, mockUserRepo)
		response, err := service.CreateEvent(request)

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if response != nil {
			t.Errorf("Expected nil response, got %v", response)
		}
		if err.Error() != "error reading users" {
			t.Errorf("Expected error message 'error reading users', got '%s'", err.Error())
		}
	})

	// Test case 4: Error when creating event
	t.Run("ErrorCreatingEvent", func(t *testing.T) {
		t.Parallel() // Enable parallel execution for this subtest

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEventRepo := mocks.NewMockEventRepository(ctrl)
		mockUserRepo := mocks.NewMockUserRepository(ctrl)

		description := "Test Description"
		location := "Test Location"
		startTime := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
		endTime := time.Now().Add(48 * time.Hour).Format(time.RFC3339)
		limit := 10

		request := &types.CreateEventRequest{
			Title:         "Test Event",
			Description:   &description,
			Location:      &location,
			StartTime:     &startTime,
			EndTime:       &endTime,
			CreatedBy:     1,
			IsPublic:      true,
			AttendeeLimit: &limit,
		}

		mockEventRepo.EXPECT().
			CreateEvent(gomock.Any()).
			Return(nil, errors.New("error creating event"))

		service := NewEventServiceImpl(mockEventRepo, mockUserRepo)
		response, err := service.CreateEvent(request)

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if response != nil {
			t.Errorf("Expected nil response, got %v", response)
		}
		if err.Error() != "error creating event" {
			t.Errorf("Expected error message 'error creating event', got '%s'", err.Error())
		}
	})
}

// Test cases for EventServiceImpl.ListEvents
func TestListEvents(t *testing.T) {
	// Test case 1: Successful listing of events for admin user
	t.Run("SuccessfulListingForAdmin", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEventRepo := mocks.NewMockEventRepository(ctrl)
		mockUserRepo := mocks.NewMockUserRepository(ctrl)

		request := types.ListEventRequest{
			Page:  1,
			Limit: 10,
		}

		user := &types.CurrentUser{
			ID:          1,
			Email:       "admin@example.com",
			RoleID:      1,
			Role:        "Admin",
			Permissions: []string{consts.PermissionFetchAllEvent},
		}

		events := []*models.Event{
			createTestEvent(1),
			createTestEvent(2),
		}

		mockEventRepo.EXPECT().
			ListEvents(gomock.Any(), gomock.Eq(10), gomock.Eq(0)).
			Return(events, 2, nil)

		service := NewEventServiceImpl(mockEventRepo, mockUserRepo)
		response, err := service.ListEvents(request, user)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if response == nil {
			t.Error("Expected response, got nil")
		}
		if response.Total != 2 {
			t.Errorf("Expected total 2, got %d", response.Total)
		}
		if response.Page != 1 {
			t.Errorf("Expected page 1, got %d", response.Page)
		}
		if response.Limit != 10 {
			t.Errorf("Expected limit 10, got %d", response.Limit)
		}
		if !reflect.DeepEqual(response.Events, events) {
			t.Errorf("Expected events %v, got %v", events, response.Events)
		}
	})

	// Test case 2: Successful listing of events for public access (no user)
	t.Run("SuccessfulListingForPublicAccess", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEventRepo := mocks.NewMockEventRepository(ctrl)
		mockUserRepo := mocks.NewMockUserRepository(ctrl)

		request := types.ListEventRequest{
			Page:  1,
			Limit: 10,
		}

		events := []*models.Event{
			createTestEvent(1),
			createTestEvent(2),
		}

		// Use a matcher to verify IsPublic is set to true
		mockEventRepo.EXPECT().
			ListEvents(gomock.Any(), gomock.Eq(10), gomock.Eq(0)).
			DoAndReturn(func(filter *types.EventFilter, limit, offset int) ([]*models.Event, int, error) {
				// Verify that IsPublic is set to true for public access
				if filter.IsPublic == nil || *filter.IsPublic != true {
					t.Error("Expected IsPublic to be true for public access")
				}
				return events, 2, nil
			})

		service := NewEventServiceImpl(mockEventRepo, mockUserRepo)
		response, err := service.ListEvents(request, nil)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if response == nil {
			t.Error("Expected response, got nil")
		}
		if response.Total != 2 {
			t.Errorf("Expected total 2, got %d", response.Total)
		}
		if response.Page != 1 {
			t.Errorf("Expected page 1, got %d", response.Page)
		}
		if response.Limit != 10 {
			t.Errorf("Expected limit 10, got %d", response.Limit)
		}
		if !reflect.DeepEqual(response.Events, events) {
			t.Errorf("Expected events %v, got %v", events, response.Events)
		}
	})

	// Test case 3: No events found
	t.Run("NoEventsFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEventRepo := mocks.NewMockEventRepository(ctrl)
		mockUserRepo := mocks.NewMockUserRepository(ctrl)

		request := types.ListEventRequest{
			Page:  1,
			Limit: 10,
		}

		user := &types.CurrentUser{
			ID:          1,
			Email:       "admin@example.com",
			RoleID:      1,
			Role:        "Admin",
			Permissions: []string{consts.PermissionFetchAllEvent},
		}

		mockEventRepo.EXPECT().
			ListEvents(gomock.Any(), gomock.Eq(10), gomock.Eq(0)).
			Return(nil, 0, errutil.ErrRecordNotFound)

		service := NewEventServiceImpl(mockEventRepo, mockUserRepo)
		response, err := service.ListEvents(request, user)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if response == nil {
			t.Error("Expected response, got nil")
		}
		if response.Total != 0 {
			t.Errorf("Expected total 0, got %d", response.Total)
		}
		if len(response.Events) != 0 {
			t.Errorf("Expected empty events, got %v", response.Events)
		}
	})

	// Test case 4: Error when listing events
	t.Run("ErrorListingEvents", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEventRepo := mocks.NewMockEventRepository(ctrl)
		mockUserRepo := mocks.NewMockUserRepository(ctrl)

		request := types.ListEventRequest{
			Page:  1,
			Limit: 10,
		}

		user := &types.CurrentUser{
			ID:          1,
			Email:       "admin@example.com",
			RoleID:      1,
			Role:        "Admin",
			Permissions: []string{consts.PermissionFetchAllEvent},
		}

		mockEventRepo.EXPECT().
			ListEvents(gomock.Any(), gomock.Eq(10), gomock.Eq(0)).
			Return(nil, 0, errors.New("error listing events"))

		service := NewEventServiceImpl(mockEventRepo, mockUserRepo)
		response, err := service.ListEvents(request, user)

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if response != nil {
			t.Errorf("Expected nil response, got %v", response)
		}
		if err.Error() != "error listing events" {
			t.Errorf("Expected error message 'error listing events', got '%s'", err.Error())
		}
	})
}

// Test cases for EventServiceImpl.ReadEventByID
func TestReadEventByID(t *testing.T) {
	// Test case 1: Successful reading of an event
	t.Run("SuccessfulReading", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEventRepo := mocks.NewMockEventRepository(ctrl)
		mockUserRepo := mocks.NewMockUserRepository(ctrl)

		expectedEvent := createTestEvent(1)

		mockEventRepo.EXPECT().
			ReadEventByID(gomock.Eq(1)).
			Return(expectedEvent, nil)

		service := NewEventServiceImpl(mockEventRepo, mockUserRepo)
		event, err := service.ReadEventByID(1)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if event == nil {
			t.Error("Expected event, got nil")
		}
		if !reflect.DeepEqual(event, expectedEvent) {
			t.Errorf("Expected event %v, got %v", expectedEvent, event)
		}
	})

	// Test case 2: Error when reading an event
	t.Run("ErrorReading", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEventRepo := mocks.NewMockEventRepository(ctrl)
		mockUserRepo := mocks.NewMockUserRepository(ctrl)

		mockEventRepo.EXPECT().
			ReadEventByID(gomock.Eq(1)).
			Return(nil, errors.New("error reading event"))

		service := NewEventServiceImpl(mockEventRepo, mockUserRepo)
		event, err := service.ReadEventByID(1)

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if event != nil {
			t.Errorf("Expected nil event, got %v", event)
		}
		if err.Error() != "error reading event" {
			t.Errorf("Expected error message 'error reading event', got '%s'", err.Error())
		}
	})
}

// Test cases for EventServiceImpl.UpdateEvent
func TestUpdateEvent(t *testing.T) {
	// Test case 1: Successful update of an event
	t.Run("SuccessfulUpdate", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEventRepo := mocks.NewMockEventRepository(ctrl)
		mockUserRepo := mocks.NewMockUserRepository(ctrl)

		description := "Updated Description"
		location := "Updated Location"
		startTime := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
		endTime := time.Now().Add(48 * time.Hour).Format(time.RFC3339)

		request := &types.UpdateEventRequest{
			ID: 1,
			CreateEventRequest: types.CreateEventRequest{
				Title:       "Updated Event",
				Description: &description,
				Location:    &location,
				StartTime:   &startTime,
				EndTime:     &endTime,
				CreatedBy:   1,
			},
		}

		existingEvent := createTestEvent(1)
		updatedEvent := request.ToEvent()
		updatedEvent.ID = 1

		mockEventRepo.EXPECT().
			ReadEventByID(gomock.Eq(1)).
			Return(existingEvent, nil)

		mockEventRepo.EXPECT().
			UpdateEvent(gomock.Any()).
			Return(updatedEvent, nil)

		service := NewEventServiceImpl(mockEventRepo, mockUserRepo)
		response, err := service.UpdateEvent(request)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if response == nil {
			t.Error("Expected response, got nil")
		}
		if response.Message != "Event updated" {
			t.Errorf("Expected message 'Event updated', got '%s'", response.Message)
		}
		if !reflect.DeepEqual(response.Event, updatedEvent) {
			t.Errorf("Expected event %v, got %v", updatedEvent, response.Event)
		}
	})

	// Test case 2: Event not found
	t.Run("EventNotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEventRepo := mocks.NewMockEventRepository(ctrl)
		mockUserRepo := mocks.NewMockUserRepository(ctrl)

		description := "Updated Description"
		location := "Updated Location"
		startTime := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
		endTime := time.Now().Add(48 * time.Hour).Format(time.RFC3339)

		request := &types.UpdateEventRequest{
			ID: 1,
			CreateEventRequest: types.CreateEventRequest{
				Title:       "Updated Event",
				Description: &description,
				Location:    &location,
				StartTime:   &startTime,
				EndTime:     &endTime,
				CreatedBy:   1,
			},
		}

		mockEventRepo.EXPECT().
			ReadEventByID(gomock.Eq(1)).
			Return(nil, nil)

		service := NewEventServiceImpl(mockEventRepo, mockUserRepo)
		response, err := service.UpdateEvent(request)

		if err != errutil.ErrRecordNotFound {
			t.Errorf("Expected error ErrRecordNotFound, got %v", err)
		}
		if response != nil {
			t.Errorf("Expected nil response, got %v", response)
		}
	})

	// Test case 3: Error when reading event
	t.Run("ErrorReadingEvent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEventRepo := mocks.NewMockEventRepository(ctrl)
		mockUserRepo := mocks.NewMockUserRepository(ctrl)

		description := "Updated Description"
		location := "Updated Location"
		startTime := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
		endTime := time.Now().Add(48 * time.Hour).Format(time.RFC3339)

		request := &types.UpdateEventRequest{
			ID: 1,
			CreateEventRequest: types.CreateEventRequest{
				Title:       "Updated Event",
				Description: &description,
				Location:    &location,
				StartTime:   &startTime,
				EndTime:     &endTime,
				CreatedBy:   1,
			},
		}

		mockEventRepo.EXPECT().
			ReadEventByID(gomock.Eq(1)).
			Return(nil, errors.New("error reading event"))

		service := NewEventServiceImpl(mockEventRepo, mockUserRepo)
		response, err := service.UpdateEvent(request)

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if response != nil {
			t.Errorf("Expected nil response, got %v", response)
		}
		if err.Error() != "error reading event" {
			t.Errorf("Expected error message 'error reading event', got '%s'", err.Error())
		}
	})

	// Test case 4: Error when updating event
	t.Run("ErrorUpdatingEvent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEventRepo := mocks.NewMockEventRepository(ctrl)
		mockUserRepo := mocks.NewMockUserRepository(ctrl)

		description := "Updated Description"
		location := "Updated Location"
		startTime := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
		endTime := time.Now().Add(48 * time.Hour).Format(time.RFC3339)

		request := &types.UpdateEventRequest{
			ID: 1,
			CreateEventRequest: types.CreateEventRequest{
				Title:       "Updated Event",
				Description: &description,
				Location:    &location,
				StartTime:   &startTime,
				EndTime:     &endTime,
				CreatedBy:   1,
			},
		}

		existingEvent := createTestEvent(1)

		mockEventRepo.EXPECT().
			ReadEventByID(gomock.Eq(1)).
			Return(existingEvent, nil)

		mockEventRepo.EXPECT().
			UpdateEvent(gomock.Any()).
			Return(nil, errors.New("error updating event"))

		service := NewEventServiceImpl(mockEventRepo, mockUserRepo)
		response, err := service.UpdateEvent(request)

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if response != nil {
			t.Errorf("Expected nil response, got %v", response)
		}
		if err.Error() != "error updating event" {
			t.Errorf("Expected error message 'error updating event', got '%s'", err.Error())
		}
	})
}

// Test cases for EventServiceImpl.DeleteEvent
func TestDeleteEvent(t *testing.T) {
	// Test case 1: Successful deletion of an event
	t.Run("SuccessfulDeletion", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEventRepo := mocks.NewMockEventRepository(ctrl)
		mockUserRepo := mocks.NewMockUserRepository(ctrl)

		mockEventRepo.EXPECT().
			DeleteEvent(gomock.Eq(1)).
			Return(nil)

		service := NewEventServiceImpl(mockEventRepo, mockUserRepo)
		response, err := service.DeleteEvent(1)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if response == nil {
			t.Error("Expected response, got nil")
		}
		if response.Message != "Event deleted" {
			t.Errorf("Expected message 'Event deleted', got '%s'", response.Message)
		}
	})

	// Test case 2: Error when deleting an event
	t.Run("ErrorDeleting", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEventRepo := mocks.NewMockEventRepository(ctrl)
		mockUserRepo := mocks.NewMockUserRepository(ctrl)

		mockEventRepo.EXPECT().
			DeleteEvent(gomock.Eq(1)).
			Return(errors.New("error deleting event"))

		service := NewEventServiceImpl(mockEventRepo, mockUserRepo)
		response, err := service.DeleteEvent(1)

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if response != nil {
			t.Errorf("Expected nil response, got %v", response)
		}
		if err.Error() != "error deleting event" {
			t.Errorf("Expected error message 'error deleting event', got '%s'", err.Error())
		}
	})
}

// Test cases for EventServiceImpl.RsvpEvent
func TestRsvpEvent(t *testing.T) {
	// Test case 1: Successful RSVP for a private event
	t.Run("SuccessfulRsvpForPrivateEvent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEventRepo := mocks.NewMockEventRepository(ctrl)
		mockUserRepo := mocks.NewMockUserRepository(ctrl)

		request := types.RsvpEventRequest{
			EventID:  1,
			UserID:   2,
			StatusID: 2, // Accepted
		}

		event := createTestEvent(1)
		event.IsPublic = false

		invitation := &models.EventAttendee{
			EventID:  1,
			UserID:   2,
			StatusID: 1, // Pending
		}

		mockEventRepo.EXPECT().
			ReadEventByID(gomock.Eq(1)).
			Return(event, nil)

		mockEventRepo.EXPECT().
			ReadEventInvitation(gomock.Eq(1), gomock.Eq(2)).
			Return(invitation, nil)

		// Verify the invitation is updated correctly
		mockEventRepo.EXPECT().
			UpsertEventInvitation(gomock.Any()).
			DoAndReturn(func(inv *models.EventAttendee) error {
				if inv.EventID != 1 || inv.UserID != 2 || inv.StatusID != 2 {
					t.Errorf("Expected invitation with EventID=1, UserID=2, StatusID=2, got %v", inv)
				}
				return nil
			})

		service := NewEventServiceImpl(mockEventRepo, mockUserRepo)
		err := service.RsvpEvent(request)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	// Test case 2: Successful RSVP for a public event
	t.Run("SuccessfulRsvpForPublicEvent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEventRepo := mocks.NewMockEventRepository(ctrl)
		mockUserRepo := mocks.NewMockUserRepository(ctrl)

		request := types.RsvpEventRequest{
			EventID:  1,
			UserID:   2,
			StatusID: 2, // Accepted
		}

		event := createTestEvent(1)
		event.IsPublic = true

		mockEventRepo.EXPECT().
			ReadEventByID(gomock.Eq(1)).
			Return(event, nil)

		mockEventRepo.EXPECT().
			GetEventAttendeesCount(gomock.Eq(1)).
			Return(5, nil)

		// Verify the invitation is created correctly
		mockEventRepo.EXPECT().
			UpsertEventInvitation(gomock.Any()).
			DoAndReturn(func(inv *models.EventAttendee) error {
				if inv.EventID != 1 || inv.UserID != 2 || inv.StatusID != 2 {
					t.Errorf("Expected invitation with EventID=1, UserID=2, StatusID=2, got %v", inv)
				}
				return nil
			})

		service := NewEventServiceImpl(mockEventRepo, mockUserRepo)
		err := service.RsvpEvent(request)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	// Test case 3: Error when event capacity is exceeded
	t.Run("ErrorEventCapacityExceeded", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEventRepo := mocks.NewMockEventRepository(ctrl)
		mockUserRepo := mocks.NewMockUserRepository(ctrl)

		request := types.RsvpEventRequest{
			EventID:  1,
			UserID:   2,
			StatusID: 2, // Accepted
		}

		event := createTestEvent(1)
		event.IsPublic = true
		limit := 5
		event.Limit = &limit

		mockEventRepo.EXPECT().
			ReadEventByID(gomock.Eq(1)).
			Return(event, nil)

		mockEventRepo.EXPECT().
			GetEventAttendeesCount(gomock.Eq(1)).
			Return(5, nil) // Already at capacity

		service := NewEventServiceImpl(mockEventRepo, mockUserRepo)
		err := service.RsvpEvent(request)

		if err != errutil.ErrEventCapacityExceeded {
			t.Errorf("Expected error ErrEventCapacityExceeded, got %v", err)
		}
	})

	// Test case 4: Error when reading event
	t.Run("ErrorReadingEvent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEventRepo := mocks.NewMockEventRepository(ctrl)
		mockUserRepo := mocks.NewMockUserRepository(ctrl)

		request := types.RsvpEventRequest{
			EventID:  1,
			UserID:   2,
			StatusID: 2,
		}

		mockEventRepo.EXPECT().
			ReadEventByID(gomock.Eq(1)).
			Return(nil, errors.New("error reading event"))

		service := NewEventServiceImpl(mockEventRepo, mockUserRepo)
		err := service.RsvpEvent(request)

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err.Error() != "error reading event" {
			t.Errorf("Expected error message 'error reading event', got '%s'", err.Error())
		}
	})

	// Test case 5: Error when invitation not found for private event
	t.Run("ErrorInvitationNotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEventRepo := mocks.NewMockEventRepository(ctrl)
		mockUserRepo := mocks.NewMockUserRepository(ctrl)

		request := types.RsvpEventRequest{
			EventID:  1,
			UserID:   2,
			StatusID: 2,
		}

		event := createTestEvent(1)
		event.IsPublic = false

		mockEventRepo.EXPECT().
			ReadEventByID(gomock.Eq(1)).
			Return(event, nil)

		mockEventRepo.EXPECT().
			ReadEventInvitation(gomock.Eq(1), gomock.Eq(2)).
			Return(nil, nil)

		service := NewEventServiceImpl(mockEventRepo, mockUserRepo)
		err := service.RsvpEvent(request)

		if err != errutil.ErrRecordNotFound {
			t.Errorf("Expected error ErrRecordNotFound, got %v", err)
		}
	})

	// Test case 6: Error when upserting invitation
	t.Run("ErrorUpsertingInvitation", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEventRepo := mocks.NewMockEventRepository(ctrl)
		mockUserRepo := mocks.NewMockUserRepository(ctrl)

		request := types.RsvpEventRequest{
			EventID:  1,
			UserID:   2,
			StatusID: 2,
		}

		event := createTestEvent(1)
		event.IsPublic = true

		mockEventRepo.EXPECT().
			ReadEventByID(gomock.Eq(1)).
			Return(event, nil)

		mockEventRepo.EXPECT().
			GetEventAttendeesCount(gomock.Eq(1)).
			Return(5, nil)

		mockEventRepo.EXPECT().
			UpsertEventInvitation(gomock.Any()).
			Return(errors.New("error upserting invitation"))

		service := NewEventServiceImpl(mockEventRepo, mockUserRepo)
		err := service.RsvpEvent(request)

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err.Error() != "error upserting invitation" {
			t.Errorf("Expected error message 'error upserting invitation', got '%s'", err.Error())
		}
	})
}

// Benchmark for EventServiceImpl.CreateEvent
func BenchmarkCreateEvent(b *testing.B) {
	description := "Benchmark Description"
	location := "Benchmark Location"
	startTime := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	endTime := time.Now().Add(48 * time.Hour).Format(time.RFC3339)
	limit := 10

	for i := 0; i < b.N; i++ {
		ctrl := gomock.NewController(b)
		mockEventRepo := mocks.NewMockEventRepository(ctrl)
		mockUserRepo := mocks.NewMockUserRepository(ctrl)

		request := &types.CreateEventRequest{
			Title:         "Benchmark Event",
			Description:   &description,
			Location:      &location,
			StartTime:     &startTime,
			EndTime:       &endTime,
			CreatedBy:     1,
			IsPublic:      true,
			AttendeeLimit: &limit,
		}

		expectedEvent := request.ToEvent()
		expectedEvent.ID = i + 1

		mockEventRepo.EXPECT().
			CreateEvent(gomock.Any()).
			Return(expectedEvent, nil)

		service := NewEventServiceImpl(mockEventRepo, mockUserRepo)
		_, err := service.CreateEvent(request)
		if err != nil {
			b.Errorf("Unexpected error: %v", err)
		}
		ctrl.Finish()
	}
}
