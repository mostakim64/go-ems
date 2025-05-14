package cmd

import (
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"github.com/vivasoft-ltd/go-ems/conn"
	"github.com/vivasoft-ltd/go-ems/controllers"
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

	// repositories
	eventRepo := db_repo.NewEventRepositoryImpl(dbClient)

	// services
	eventSvc := services.NewEventServiceImpl(eventRepo)

	// controllers
	eventCtrl := controllers.NewEventController(eventSvc)

	// Server
	var echo_ = echo.New()
	var Routes = routes.New(echo_, eventCtrl)
	var Server = server.New(echo_)

	// Spooling
	Routes.Init()
	Server.Start()
}
