package domain

import (
	"github.com/vivasoft-ltd/go-ems/models"
	"github.com/vivasoft-ltd/go-ems/types"
)

type (
	MailService interface {
		SendEmail(reqData types.EmailPayload) error
		SendInvitationEmail(userIds []int, event *models.Event) error
		EnqueueEventReminderEmailNotification(event *models.Event) error
	}

	MailRepository interface {
		SendEmail(reqData *types.EmailPayload) error
	}
)
