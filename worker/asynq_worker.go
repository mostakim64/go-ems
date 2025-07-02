package worker

import (
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/vivasoft-ltd/go-ems/config"
	"github.com/vivasoft-ltd/go-ems/types"
)

func StartAsynqWorker(mux *asynq.ServeMux) {
	worker := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     config.Asynq().RedisAddr,
			DB:       config.Asynq().DB,
			Password: config.Asynq().Pass,
		},
		asynq.Config{
			Concurrency: config.Asynq().Concurrency,
			Queues: map[string]int{
				config.Asynq().Queue: 1,
			},
			RetryDelayFunc: func(numOfRetry int, e error, t *asynq.Task) time.Duration {
				switch t.Type() {
				case types.AsynqTaskTypeInvitationEmail.String():
					return config.Asynq().EmailInvitationTaskRetryDelay * time.Second
				case types.AsynqTaskTypeEventReminder.String():
					return config.Asynq().EventReminderTaskRetryDelay * time.Second
				case types.AsynqTaskTypeEventReminderEmail.String():
					return config.Asynq().EventReminderEmailTaskRetryDelay * time.Second
				default:
					return asynq.DefaultRetryDelayFunc(numOfRetry, e, t)
				}
			},
		},
	)

	if err := worker.Run(mux); err != nil {
		panic(fmt.Sprintf("could not run worker: %v", err))
	}
}
