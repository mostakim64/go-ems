package cmd

import (
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"github.com/vivasoft-ltd/go-ems/config"
	"github.com/vivasoft-ltd/go-ems/conn"
	"github.com/vivasoft-ltd/go-ems/controllers"
	"github.com/vivasoft-ltd/go-ems/middlewares"
	mail_repo "github.com/vivasoft-ltd/go-ems/repositories/mail"

	asynq_repo "github.com/vivasoft-ltd/go-ems/repositories/asynq"
	db_repo "github.com/vivasoft-ltd/go-ems/repositories/db"
	"github.com/vivasoft-ltd/go-ems/routes"
	"github.com/vivasoft-ltd/go-ems/server"
	"github.com/vivasoft-ltd/go-ems/services"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: serve,
}

func serve(cmd *cobra.Command, args []string) {
	// clients
	dbClient := conn.Db()
	redisClient := conn.Redis()
	emailClient := conn.EmailClient()
	asynqClient := conn.Asynq()
	asynqInspector := conn.AsynqInspector()

	// repositories
	dbRepo := db_repo.NewRepository(dbClient)
	asynqRepo := asynq_repo.NewRepository(config.Asynq(), asynqClient, asynqInspector)
	mailRepo := mail_repo.NewRepository(emailClient, config.Email())

	// services
	redisSvc := services.NewRedisService(redisClient)
	eventSvc := services.NewEventServiceImpl(dbRepo, dbRepo)
	userSvc := services.NewUserServiceImpl(redisSvc, dbRepo)
	tokenSvc := services.NewTokenServiceImpl(redisSvc)
	authSvc := services.NewAuthServiceImpl(userSvc, tokenSvc)
	mailSvc := services.NewMailService(dbRepo, dbRepo, mailRepo)
	asynqSvc := services.NewAsynqService(config.Asynq(), asynqRepo, dbRepo, dbRepo)

	// controllers
	eventCtrl := controllers.NewEventController(eventSvc, mailSvc, asynqSvc)
	userCtrl := controllers.NewUserController(userSvc)
	authCtrl := controllers.NewAuthController(authSvc)

	// middlewares
	authMiddleware := middlewares.NewAuthMiddleware(authSvc, userSvc)

	// Server
	var echo_ = echo.New()
	var Routes = routes.New(echo_, eventCtrl, userCtrl, authCtrl, authMiddleware)
	var Server = server.New(echo_)

	// Spooling
	Routes.Init()

	// Stopping running workers
	Server.Start()
}
