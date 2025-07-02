package services

import (
	"fmt"
	"time"

	"github.com/vivasoft-ltd/go-ems/consts"
	"github.com/vivasoft-ltd/go-ems/domain"
	"github.com/vivasoft-ltd/go-ems/models"
	"github.com/vivasoft-ltd/go-ems/types"
	"github.com/vivasoft-ltd/go-ems/utils/errutil"
	"github.com/vivasoft-ltd/go-ems/worker"
	"github.com/vivasoft-ltd/golang-course-utils/logger"
)

type Mail struct {
	userRepo   domain.UserRepository
	eventRepo  domain.EventRepository
	mailRepo   domain.MailRepository
	workerPool *worker.Pool
}

func NewMailService(userRepo domain.UserRepository, eventRepo domain.EventRepository, mailRepo domain.MailRepository) *Mail {
	return &Mail{
		userRepo:  userRepo,
		eventRepo: eventRepo,
		mailRepo:  mailRepo,
	}
}

func (m *Mail) SendEmail(reqData types.EmailPayload) error {
	err := m.mailRepo.SendEmail(&reqData)
	if err != nil {
		logger.Error(fmt.Sprintf("err: []%v occurred while sending email to: %s", err, reqData.MailTo))
		return err
	}
	return nil
}

// func (m *Mail) SendInvitationEmail(userIds []int, event *models.Event) error {
// 	users, err := m.userRepo.ReadUsers(userIds)
// 	if err != nil {
// 		return err
// 	}
// 	for _, user := range users {
// 		emailPayload := types.EmailPayload{
// 			MailTo:  user.Email,
// 			Subject: "Invitation to Event: " + event.Title,
// 			Body: map[string]interface{}{
// 				"event":     event,
// 				"time":      event.StartTime.Format("2006-01-02 15:04:05"),
// 				"rsvp_link": fmt.Sprintf("http://127.0.0.1:8080/v1/events/%d/rsvp", event.ID),
// 			},
// 		}

// 		m.SendEmail(emailPayload)
// 	}

// 	return nil
// }

func (m *Mail) SendInvitationEmail(userIds []int, event *models.Event) error {
	users, err := m.userRepo.ReadUsers(userIds)
	if err != nil {
		return err
	}

	for _, user := range users {
		emailPayload := types.EmailPayload{
			MailTo:  user.Email,
			Subject: "Invitation to Event: " + event.Title,
			Body: map[string]interface{}{
				"event":     event,
				"rsvp_link": fmt.Sprintf("http://127.0.0.1:8080/v1/events/%d/rsvp", event.ID),
			},
		}

		// Add the email sending task to the worker pool
		task := worker.NewTask(func() error {
			return m.SendEmail(emailPayload)
		}, func(err error) {
			logger.Error("Failed to send email: ", err, " to user: ", user.Email)
		}, 3)

		m.workerPool.AddTask(task)
	}

	// wp.StopAfterTaskCompleted(len(users))

	return nil
}

func (m *Mail) EnqueueEventReminderEmailNotification(event *models.Event) error {
	eventStartTime := event.StartTime
	reminderTime := eventStartTime.Add(-consts.EventReminderInterval)
	now := time.Now()

	if reminderTime.Before(now) {
		logger.Info("Reminder time is in the past, skipping event reminder email for event: ", event.Title)
		return errutil.ErrEventReminderEmailNotEnqueued
	}

	logger.Info("Enqueing event reminder email for event: ", event.Title)
	logger.Info(fmt.Sprintf("Reminder trigger time: %v, event start time: %v", reminderTime.Format(time.RFC3339), eventStartTime.Format(time.RFC3339)))

	timeLeftToSendReminderEmail := reminderTime.Sub(now)

	// Initialize the scheduler
	scheduler := worker.NewScheduler(timeLeftToSendReminderEmail)
	// pass callback/callable function to the scheduler
	scheduler.Start(func() { m.sendEventReminderEmail(event) })

	// Let the scheduler run for some time
	time.Sleep(time.Duration(timeLeftToSendReminderEmail))
	scheduler.Stop()

	return nil
}

func (m *Mail) sendEventReminderEmail(event *models.Event) {
	eventAttendees, err := m.eventRepo.GetAcceptedEventAttendees(event.ID)
	if err != nil {
		if err == errutil.ErrUserNotFound {
			logger.Error(fmt.Sprintf("SKIPPING: No accepted event attendees found for event: %s", event.Title))
			return
		}

		logger.Error(fmt.Sprintf("Failed to get accepted event attendees for event: %s , err: %v", event.Title, err))
		return
	}

	fmt.Printf("Sending event reminder email to %d Attendees who have accepted the invitation", len(eventAttendees))

	for _, eventAttendee := range eventAttendees {
		user := eventAttendee.User
		emailPayload := types.EmailPayload{
			MailTo:  user.Email,
			Subject: "Event Reminder: " + event.Title,
			Body: map[string]interface{}{
				"event_title": event.Title,
				"start_time":  event.StartTime,
				"join_link":   "https://www.go-ems.com/join?_C=dQw4w9WgXcQ",
			},
		}

		task := worker.NewTask(func() error {
			return m.SendEmail(emailPayload)
		}, func(err error) {
			logger.Error("Failed to send reminder email: ", err, " to user: ", user.Email)
		}, 0)

		m.workerPool.AddTask(task)
	}
}
