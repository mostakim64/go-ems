package domain

import "github.com/vivasoft-ltd/go-ems/models"

type (
	MailService interface {
		SendInvitationEmail(userIds []int, event *models.Event) error
	}
)
