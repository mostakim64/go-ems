package routes

import (
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/vivasoft-ltd/go-ems/consts"
	"github.com/vivasoft-ltd/go-ems/controllers"
	m "github.com/vivasoft-ltd/go-ems/middlewares"
)

type Routes struct {
	echo           *echo.Echo
	eventCtrl      *controllers.EventController
	userCtrl       *controllers.UserController
	authCtrl       *controllers.AuthController
	authMiddleware *m.AuthMiddleware
}

func New(e *echo.Echo, eventCtrl *controllers.EventController, userCtrl *controllers.UserController, authCtrl *controllers.AuthController, authMiddleware *m.AuthMiddleware) *Routes {
	return &Routes{
		echo:           e,
		eventCtrl:      eventCtrl,
		userCtrl:       userCtrl,
		authCtrl:       authCtrl,
		authMiddleware: authMiddleware,
	}
}

func (r *Routes) Init() {
	e := r.echo
	m.Init(e)
	// APM routes
	e.GET("/metrics", echoprometheus.NewHandler())

	g := e.Group("/v1")

	g.POST("/events", r.eventCtrl.CreateEvent, r.authMiddleware.Authenticate(consts.PermissionEventCreate))
	g.GET("/events", r.eventCtrl.ListEvents, r.authMiddleware.Authenticate(consts.PermissionEventList))
	g.GET("/events/public", r.eventCtrl.ListPublicEvents)
	g.GET("/events/:id", r.eventCtrl.ReadEventByID, r.authMiddleware.Authenticate(consts.PermissionEventFetch))
	g.PUT("/events/:id", r.eventCtrl.UpdateEvent, r.authMiddleware.Authenticate(consts.PermissionEventUpdate))
	g.DELETE("/events/:id", r.eventCtrl.DeleteEvent, r.authMiddleware.Authenticate(consts.PermissionEventDelete))
	g.POST("/events/:id/rsvp", r.eventCtrl.Rsvp, r.authMiddleware.Authenticate(""))

	users := g.Group("/users")
	users.POST("/signup", r.userCtrl.Signup)
	users.GET("/profile", r.userCtrl.Profile, r.authMiddleware.Authenticate(""))
	users.POST("", r.userCtrl.CreateUser, r.authMiddleware.Authenticate(consts.PermissionUserCreate))
	users.GET("", r.userCtrl.ListUsers, r.authMiddleware.Authenticate(consts.PermissionUserList))
	users.GET("/:id", r.userCtrl.ReadUser, r.authMiddleware.Authenticate(consts.PermissionUserFetch))
	users.PUT("/:id", r.userCtrl.UpdateUser, r.authMiddleware.Authenticate(consts.PermissionUserUpdate))
	users.DELETE("/:id", r.userCtrl.DeleteUser, r.authMiddleware.Authenticate(consts.PermissionUserDelete))
	users.GET("/attendees", r.userCtrl.ListAttendees, r.authMiddleware.Authenticate(consts.PermissionListAttendee))

	auth := g.Group("/auth")
	auth.POST("/login", r.authCtrl.Login)
	auth.POST("/logout", r.authCtrl.Logout, r.authMiddleware.Authenticate(""))

}
