package middlewares

import (
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	m "github.com/labstack/echo/v4/middleware"
)

const (
	promSubsystemName = "vivasoft"
)

func Init(e *echo.Echo) {
	var (
		metricsPath string = "/metrics"
	)

	e.Pre(m.RemoveTrailingSlash())
	e.Use(m.LoggerWithConfig(m.LoggerConfig{
		Format:           `${time_custom} ${remote_ip} ${host} ${method} ${uri} ${status} ${latency_human} ${bytes_in} ${bytes_out} "${user_agent}"` + "\n",
		CustomTimeFormat: "2006-01-02T15:04:05.00",
	}))
	e.Use(m.CORS())
	e.Use(m.Secure())
	e.Use(m.Recover())
	e.Use(m.GzipWithConfig(m.GzipConfig{
		Skipper: func(context echo.Context) bool {
			return context.Request().URL.Path == metricsPath
		},
		Level: 5,
	}))
	e.Use(echoprometheus.NewMiddleware(promSubsystemName)) // adds middleware to gather metrics
}
