package db

import (
	"errors"
	"fmt"
	"github.com/vivasoft-ltd/go-ems/models"
	"github.com/vivasoft-ltd/go-ems/types"
	"github.com/vivasoft-ltd/go-ems/utils/errutil"
	"github.com/vivasoft-ltd/golang-course-utils/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (repo *Repository) CreateEvent(event *models.Event) (*models.Event, error) {
	qry := repo.client.Create(event)
	if qry.Error != nil {
		logger.Error(fmt.Errorf("error creating event: %w", qry.Error))
		return nil, qry.Error
	}

	return event, nil
}

func (repo *Repository) ListEvents(filter *types.EventFilter, Limit, Offset int) ([]*models.Event, int, error) {
	query := repo.client.Model(&models.Event{})
	fmt.Printf("+ %+v", filter)
	if filter != nil {
		if filter.IsPublic != nil {
			query = query.Where("is_public = ?", filter.IsPublic)
		}
		if filter.CreatedBy != nil {
			query = query.Where("created_by = ?", filter.CreatedBy)
		}
		if filter.Attendee != nil {
			query = query.Where("(is_public = ? and end_time < curdate()) OR id IN (SELECT event_id FROM event_attendees WHERE user_id = ?)", true, filter.Attendee)
		}
	}
	var events []*models.Event
	var count int64
	if err := repo.client.Model(&models.Event{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	result := query.Offset(Offset).Limit(Limit).Find(&events)
	if result.RowsAffected == 0 {
		logger.Error("no events found")
		return nil, 0, errutil.ErrRecordNotFound
	}
	if result.Error != nil {
		logger.Error(fmt.Errorf("error listing events: %w", result.Error))
		return nil, 0, result.Error
	}

	return events, int(count), nil
}

func (repo *Repository) ReadEventByID(id int) (*models.Event, error) {
	var event models.Event
	qry := repo.client.Preload("Attendees").First(&event, id)
	if errors.Is(qry.Error, gorm.ErrRecordNotFound) {
		logger.Error(fmt.Errorf("event with ID %d not found", id))
		return nil, errutil.ErrRecordNotFound
	}
	if qry.Error != nil {
		logger.Error(fmt.Errorf("error getting event by ID: %w", qry.Error))
		return nil, qry.Error
	}

	return &event, nil
}

func (repo *Repository) UpdateEvent(event *models.Event) (*models.Event, error) {
	qry := repo.client.Where("id = ?", event.ID).Updates(event)
	if errors.Is(qry.Error, gorm.ErrRecordNotFound) {
		logger.Error(fmt.Errorf("no event found with ID %d", event.ID))
		return nil, errutil.ErrRecordNotFound
	}
	if qry.Error != nil {
		logger.Error(fmt.Errorf("error updating event: %w", qry.Error))
		return nil, qry.Error
	}
	return event, nil
}

func (repo *Repository) DeleteEvent(id int) error {
	qry := repo.client.Where("id = ?", id).Delete(&models.Event{})
	if qry.RowsAffected == 0 {
		logger.Error(fmt.Errorf("no event found with ID %d", id))
		return errutil.ErrRecordNotFound
	}
	if qry.Error != nil {
		logger.Error(fmt.Errorf("error deleting event: %w", qry.Error))
		return qry.Error
	}
	return nil
}
func (repo *Repository) ReadEventInvitation(event int, userID int) (*models.EventAttendee, error) {
	var invitation models.EventAttendee
	if err := repo.client.Where("event_id = ? AND user_id = ?", event, userID).First(&invitation).Error; err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (repo *Repository) UpsertEventInvitation(invitation *models.EventAttendee) error {
	qry := repo.client.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "event_id"}, {Name: "user_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"status_id": invitation.StatusID}),
	}).Create(invitation)

	if qry.Error != nil {
		logger.Error(fmt.Errorf("error upserting event invitation: %w", qry.Error))
		return qry.Error
	}
	return nil
}
func (repo *Repository) GetEventAttendeesCount(eventID int) (int, error) {
	var count int64
	if err := repo.client.Model(&models.EventAttendee{}).Where("event_id = ?", eventID).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
