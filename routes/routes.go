package routes

import (
	"github.com/vivasoft-ltd/go-ems/controllers"
	m "github.com/vivasoft-ltd/go-ems/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Routes struct {
	echo      *echo.Echo
	eventCtrl *controllers.EventController
}

func New(e *echo.Echo, eventCtrl *controllers.EventController) *Routes {
	return &Routes{
		echo:      e,
		eventCtrl: eventCtrl,
	}
}

func (r *Routes) Init() {
	e := r.echo
	m.Init(e)

	// APM routes
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	g := e.Group("/v1")

	g.POST("/events", r.eventCtrl.CreateEvent)
	g.GET("/events", r.eventCtrl.ListEvents)
	g.GET("/events/:id", r.eventCtrl.ReadEventByID)
	g.PUT("/events/:id", r.eventCtrl.UpdateEvent)
	g.DELETE("/events/:id", r.eventCtrl.DeleteEvent)
}
