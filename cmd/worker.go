package cmd

import (
	asynq_ "github.com/hibiken/asynq"
	"github.com/spf13/cobra"
	"github.com/vivasoft-ltd/go-ems/config"
	"github.com/vivasoft-ltd/go-ems/conn"
	"github.com/vivasoft-ltd/go-ems/controllers"
	asynq_repo "github.com/vivasoft-ltd/go-ems/repositories/asynq"
	db_repo "github.com/vivasoft-ltd/go-ems/repositories/db"
	mail_repo "github.com/vivasoft-ltd/go-ems/repositories/mail"
	"github.com/vivasoft-ltd/go-ems/services"
	"github.com/vivasoft-ltd/go-ems/types"
	"github.com/vivasoft-ltd/go-ems/worker"
)

var workerCmd = &cobra.Command{
	Use: "worker",
	Run: runWorker,
}

func runWorker(cmd *cobra.Command, args []string) {
	// clients
	dbClient := conn.Db()
	emailClient := conn.EmailClient()
	asynqClient := conn.Asynq()
	asynqInspector := conn.AsynqInspector()

	// repositories
	dbRepo := db_repo.NewRepository(dbClient)
	asynqRepo := asynq_repo.NewRepository(config.Asynq(), asynqClient, asynqInspector)
	mailRepo := mail_repo.NewRepository(emailClient, config.Email())

	// services
	mailSvc := services.NewMailService(dbRepo, dbRepo, mailRepo)
	asynqSvc := services.NewAsynqService(config.Asynq(), asynqRepo, dbRepo, dbRepo)

	// controllers
	asynqCtrl := controllers.NewAsynqController(mailSvc, asynqSvc)

	mux := asynq_.NewServeMux()

	mux.HandleFunc(types.AsynqTaskTypeInvitationEmail.String(), asynqCtrl.ProcessInvitationEmailTask)
	mux.HandleFunc(types.AsynqTaskTypeEventReminder.String(), asynqCtrl.ProcessEventReminderTask)
	mux.HandleFunc(types.AsynqTaskTypeEventReminderEmail.String(), asynqCtrl.ProcessEventReminderTask)
	// Start the Asynq worker
	worker.StartAsynqWorker(mux)

}
