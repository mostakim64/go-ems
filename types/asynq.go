package types

import (
	"time"
)

type (
	AsynqOption struct {
		TaskID           string
		Retry            int
		Queue            string
		RetentionHours   time.Duration
		DelaySeconds     time.Duration
		UniqueTTLSeconds time.Duration
	}

	AsynqTaskType string
)

func (t AsynqTaskType) String() string {
	return string(t)
}

const (
	AsynqTaskTypeInvitationEmail    AsynqTaskType = "go:ems:invitation_email"
	AsynqTaskTypeEventReminder      AsynqTaskType = "go:ems:event_reminder"
	AsynqTaskTypeEventReminderEmail AsynqTaskType = "go:ems:event_reminder_email"
)
