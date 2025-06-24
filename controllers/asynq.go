package controllers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/vivasoft-ltd/go-ems/domain"
	"github.com/vivasoft-ltd/go-ems/models"
	"github.com/vivasoft-ltd/go-ems/types"
	"github.com/vivasoft-ltd/golang-course-utils/logger"
)

type AsynqController struct {
	mailSvc  domain.MailService
	asynqSvc domain.AsynqService
}

func NewAsynqController(mailSvc domain.MailService, asynqSvc domain.AsynqService) *AsynqController {
	return &AsynqController{
		mailSvc:  mailSvc,
		asynqSvc: asynqSvc,
	}
}

func (ac *AsynqController) ProcessInvitationEmailTask(ctx context.Context, t *asynq.Task) (err error) {
	logger.Info(fmt.Sprintf("Received task event [%s] with ID [%s]", t.Type(), t.ResultWriter().TaskID()))
	var payload types.EmailPayload

	if err = json.Unmarshal(t.Payload(), &payload); err != nil {
		logger.Error(err)
		return
	}

	if err = ac.mailSvc.SendEmail(payload); err != nil {
		logger.Error(fmt.Sprintf("err: [%v] occurred while sending email to: %s", err, payload.MailTo))
		return err
	}
	t.ResultWriter().Write([]byte(fmt.Sprintf("Email sent successfully to %s", payload.MailTo)))
	return
}

func (ac *AsynqController) ProcessEventReminderTask(ctx context.Context, t *asynq.Task) (err error) {
	logger.Info(fmt.Sprintf("Received task event [%s] with ID [%s]", t.Type(), t.ResultWriter().TaskID()))
	var payload models.Event

	if err = json.Unmarshal(t.Payload(), &payload); err != nil {
		logger.Error(err)
		return
	}

	if err = ac.asynqSvc.CreateEventReminderEmailTasks(&payload); err != nil {
		logger.Error(fmt.Sprintf("err: [%v] occurred while creating event reminder email tasks for event: %s", err, payload.Title))
		return err
	}
	t.ResultWriter().Write([]byte(fmt.Sprintf("Event reminder email tasks created successfully for event: %s", payload.Title)))

	return
}

func (ac *AsynqController) ProcessEventReminderEmailTask(ctx context.Context, t *asynq.Task) (err error) {
	logger.Info(fmt.Sprintf("Received task event [%s] with ID [%s]", t.Type(), t.ResultWriter().TaskID()))
	var payload types.EmailPayload

	if err = json.Unmarshal(t.Payload(), &payload); err != nil {
		logger.Error(err)
		return
	}

	if err = ac.mailSvc.SendEmail(payload); err != nil {
		logger.Error(fmt.Sprintf("err: [%v] occurred while sending email to: %s", err, payload.MailTo))
		return err
	}
	t.ResultWriter().Write([]byte(fmt.Sprintf("Event reminder email sent successfully to %s", payload.MailTo)))
	return
}
