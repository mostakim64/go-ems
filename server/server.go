package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/vivasoft-ltd/go-ems/config"
	"github.com/vivasoft-ltd/golang-course-utils/logger"
)

type Server struct {
	echo *echo.Echo
}

func (s *Server) Start() {
	e := s.echo
	// Start server in a goroutine
	go func() {
		if err := e.Start(":" + config.App().Port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal("Server shutdown with error: %v", err)
	}

	logger.Info("Server exited gracefully")
}

func New(echo *echo.Echo) *Server {
	return &Server{echo: echo}
}
