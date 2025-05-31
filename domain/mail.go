package domain

type (
	MailService interface {
		SendInvitationEmail(userIds []int) error
	}
)
