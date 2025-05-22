package db

import (
	"errors"
	"fmt"
	"github.com/vivasoft-ltd/go-ems/models"
	"github.com/vivasoft-ltd/go-ems/utils/errutil"
	"github.com/vivasoft-ltd/golang-course-utils/logger"
	"gorm.io/gorm"
)

func (repo *Repository) CreateEvent(event *models.Event) (*models.Event, error) {
	qry := repo.client.Create(event)
	if qry.Error != nil {
		logger.Error(fmt.Errorf("error creating event: %w", qry.Error))
		return nil, qry.Error
	}

	return event, nil
}

func (repo *Repository) ListEvents() ([]*models.Event, error) {
	var events []*models.Event
	qry := repo.client.Find(&events)
	if qry.RowsAffected == 0 {
		logger.Error("no events found")
		return nil, errutil.ErrRecordNotFound
	}
	if qry.Error != nil {
		logger.Error(fmt.Errorf("error listing events: %w", qry.Error))
		return nil, qry.Error
	}

	return events, nil
}

func (repo *Repository) ReadEventByID(id int) (*models.Event, error) {
	var event models.Event
	qry := repo.client.First(&event, id)
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
