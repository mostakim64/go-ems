package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/vivasoft-ltd/go-ems/config"
	"github.com/vivasoft-ltd/go-ems/consts"
	"github.com/vivasoft-ltd/go-ems/domain"
	"github.com/vivasoft-ltd/go-ems/models"
	"github.com/vivasoft-ltd/go-ems/types"
	"github.com/vivasoft-ltd/go-ems/utils/errutil"
	"github.com/vivasoft-ltd/golang-course-utils/logger"
)

type AsynqService struct {
	config    *config.AsynqConfig
	asynqRepo domain.AsynqRepository
	userRepo  domain.UserRepository
	eventRepo domain.EventRepository
}

func NewAsynqService(
	config *config.AsynqConfig,
	asynqRepo domain.AsynqRepository,
	userRepo domain.UserRepository,
	eventRepo domain.EventRepository,
) *AsynqService {
	return &AsynqService{
		config:    config,
		asynqRepo: asynqRepo,
		userRepo:  userRepo,
		eventRepo: eventRepo,
	}
}

func (svc *AsynqService) CreateEmailInvitationTasks(userIds []int, event *models.Event) error {
	users, err := svc.userRepo.ReadUsers(userIds)
	if errors.Is(err, errutil.ErrUserNotFound) {
		logger.Error(fmt.Sprintf("SKIPPING: No accepted event attendees found for event: %s", event.Title))
		return nil
	}

	if err != nil {
		logger.Error(fmt.Sprintf("err: [%v] occurred while reading users for creating email invitation tasks of event id: %v", err, event.ID))
		return err
	}

	for _, user := range users {
		task, err := svc.createEmailInvitationTask(user, event)
		if err != nil {
			logger.Error(fmt.Sprintf("err: [%v] occurred while creating email invitation task for user: %v", err, user.Email))
			return err
		}

		taskID := fmt.Sprintf("%s_user:%d_event:%d", types.AsynqTaskTypeInvitationEmail, user.ID, event.ID)
		customOpts := &types.AsynqOption{
			Queue:        svc.config.Queue,
			TaskID:       taskID,
			DelaySeconds: svc.config.EmailInvitationTaskDelay,
			Retry:        svc.config.EmailInvitationTaskRetryCount,
		}
		_, err = svc.enqueueTask(task, customOpts)
		if err != nil {
			return err
		}
		logger.Info(fmt.Sprintf("enqueued email invitation task for user [%s] successfully", user.Email))
	}

	return nil
}

func (svc *AsynqService) CreateEventReminderTask(event *models.Event) error {
	eventStartTime := event.StartTime.UTC()
	reminderTime := eventStartTime.Add(-consts.EventReminderInterval)
	now := time.Now().UTC()

	if reminderTime.Before(now) {
		logger.Info("Reminder time is in the past, skipping event reminder email for event: ", event.Title)
		return errutil.ErrEventReminderEmailNotEnqueued
	}

	logger.Info("Enqueuing event reminder email for event: ", event.Title)
	logger.Info(fmt.Sprintf("Reminder trigger time: %v, event start time: %v", reminderTime.Format(time.RFC3339), eventStartTime.Format(time.RFC3339)))

	timeLeftToSendReminderEmail := time.Duration(reminderTime.Sub(now).Seconds())

	taskID := fmt.Sprintf("%s_event:%d", types.AsynqTaskTypeEventReminder, event.ID)
	customOpts := &types.AsynqOption{
		Queue:        svc.config.Queue,
		TaskID:       taskID,
		DelaySeconds: timeLeftToSendReminderEmail,
		Retry:        svc.config.EventReminderTaskRetryCount,
	}

	task, err := svc.asynqRepo.CreateTask(types.AsynqTaskTypeEventReminder, event)
	if err != nil {
		logger.Error(fmt.Sprintf("error: [%v] occurred while creating event reminder email task for event id: %d", err, event.ID))
		return err
	}

	_, err = svc.enqueueTask(task, customOpts)
	if err != nil {
		logger.Error(fmt.Sprintf("error: [%v] occurred while enqueuing event reminder email task for event id: %d", err, event.ID))
		return err
	}
	return nil
}

func (svc *AsynqService) CreateEventReminderEmailTasks(event *models.Event) error {
	eventAttendees, err := svc.eventRepo.GetAcceptedEventAttendees(event.ID)
	if errors.Is(err, errutil.ErrUserNotFound) {
		logger.Error(fmt.Sprintf("SKIPPING: No accepted event attendees found for event: %s", event.Title))
		return nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get accepted event attendees for event: %s , err: %v", event.Title, err))
		return err
	}

	for _, attendee := range eventAttendees {
		task, err := svc.createEventReminderEmailTask(attendee.User, event)
		if err != nil {
			logger.Error(fmt.Sprintf("err: [%v] occurred while creating event reminder email task for user: %v", err, attendee.User.Email))
			return err
		}
		taskID := fmt.Sprintf("%s_user:%d_event:%d", types.AsynqTaskTypeEventReminderEmail, attendee.User.ID, event.ID)
		customOpts := &types.AsynqOption{
			Queue:        svc.config.Queue,
			TaskID:       taskID,
			DelaySeconds: svc.config.EventReminderEmailTaskDelay,
			Retry:        svc.config.EventReminderEmailTaskRetryCount,
		}
		_, err = svc.enqueueTask(task, customOpts)
		if err != nil {
			return err
		}
		logger.Info(fmt.Sprintf("enqueued event reminder email task for user [%s] successfully", attendee.User.Email))
	}
	return nil
}

func (svc *AsynqService) createEmailInvitationTask(user models.User, event *models.Event) (*asynq.Task, error) {
	emailPayload := types.EmailPayload{
		MailTo:  user.Email,
		Subject: "Invitation to Event: " + event.Title,
		Body: map[string]interface{}{
			"event":     event,
			"rsvp_link": fmt.Sprintf("http://127.0.0.1:8080/v1/events/%d/rsvp", event.ID),
		},
	}

	return svc.asynqRepo.CreateTask(types.AsynqTaskTypeInvitationEmail, emailPayload)
}

func (svc *AsynqService) createEventReminderEmailTask(user models.User, event *models.Event) (*asynq.Task, error) {
	emailPayload := types.EmailPayload{
		MailTo:  user.Email,
		Subject: "Event Reminder: " + event.Title,
		Body: map[string]interface{}{
			"event_title": event.Title,
			"start_time":  event.StartTime,
			"join_link":   "https://www.go-ems.com/join?_C=dQw4w9WgXcQ",
		},
	}
	return svc.asynqRepo.CreateTask(types.AsynqTaskTypeEventReminderEmail, emailPayload)
}

func (svc *AsynqService) enqueueTask(task *asynq.Task, customOpts *types.AsynqOption) (taskID string, err error) {
	err = svc.asynqRepo.DequeueTask(customOpts.TaskID) // Ensure no duplicate tasks
	if err != nil && !errors.Is(err, asynq.ErrTaskNotFound) {
		logger.Error(fmt.Sprintf("error: [%v] occurred while dequeuing task with ID: %s", err, customOpts.TaskID))
	}

	taskID, err = svc.asynqRepo.EnqueueTask(task, customOpts)
	if errors.Is(err, asynq.ErrDuplicateTask) {
		logger.Warn(fmt.Sprintf("skipped: duplicate task for taskID: [%s]", customOpts.TaskID))
		err = nil // No error for duplicate tasks, just skip
		return
	}
	if err != nil {
		logger.Error(fmt.Sprintf("error: [%v] occurred while enqueuing task with ID: %s", err, customOpts.TaskID))
		return
	}

	logger.Info(fmt.Sprintf("enqueued task [%s] successfully", taskID))
	return taskID, nil
}
