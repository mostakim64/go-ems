package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vivasoft-ltd/go-ems/config"
	"github.com/vivasoft-ltd/go-ems/conn"
	"github.com/vivasoft-ltd/golang-course-utils/logger"
)

var (
	RootCmd = &cobra.Command{
		Use: "app",
	}
)

func init() {
	RootCmd.AddCommand(serveCmd)
}

// Execute executes the root command
func Execute() {
	// load config
	config.LoadConfig()

	// Initialize logger
	initLogger()

	conn.ConnectDb()
	conn.ConnectRedis()
	conn.InitAsynqClient()
	conn.InitAsyncInspector()

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initLogger() {
	fmt.Println("Initializing logger...")
	fmt.Println("Logger file path:", config.Logger().FilePath)
	logger.SetFileLogger(config.Logger().FilePath)
}
