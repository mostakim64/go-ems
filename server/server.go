package server

import (
	"github.com/labstack/echo/v4"
	"github.com/vivasoft-ltd/go-ems/config"
)

type Server struct {
	echo *echo.Echo
}

func (s *Server) Start() {
	e := s.echo
	e.Logger.Fatal(e.Start(":" + config.App().Port))
}

func New(echo *echo.Echo) *Server {
	return &Server{echo: echo}
}
