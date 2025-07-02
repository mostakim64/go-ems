package domain

import (
	"github.com/hibiken/asynq"
	"github.com/vivasoft-ltd/go-ems/models"
	"github.com/vivasoft-ltd/go-ems/types"
)

type (
	AsynqRepository interface {
		CreateTask(event types.AsynqTaskType, payload interface{}) (*asynq.Task, error)
		EnqueueTask(task *asynq.Task, customOpts *types.AsynqOption) (string, error)
		DequeueTask(taskID string) error
	}

	AsynqService interface {
		CreateEmailInvitationTasks(userIds []int, event *models.Event) error
		CreateEventReminderTask(event *models.Event) error
		CreateEventReminderEmailTasks(event *models.Event) error
	}
)
