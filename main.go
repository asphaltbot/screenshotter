package main

import (
	"fmt"
	"github.com/asphaltbot/screenshotter/routes"
	"github.com/asphaltbot/screenshotter/util"
	"github.com/getsentry/sentry-go"
	sentryiris "github.com/getsentry/sentry-go/iris"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"net/http"
)

var app *iris.Application

func main() {
	config := util.ReadConfig()

	if config.SentryDsn != "" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn: config.SentryDsn,
		}); err != nil {
			fmt.Printf("sentry initialisation error: %v\n", err)
		}
	}

	app = iris.New()

	if config.CORS.Enable {
		corsMiddleware := cors.New(cors.Options{
			AllowedOrigins:   config.CORS.AllowedOrigins,
			AllowedMethods:   []string{"GET", "POST"},
			AllowCredentials: false,
		})

		app.Use(corsMiddleware)
	}

	app.Use(recover.New())
	app.Use(sentryiris.New(sentryiris.Options{}))

	routes.RegisterScreenshotRoute(app)

	if config.RequestLogging.Enable {
		requestLogger := logger.New(logger.Config{
			Status: config.RequestLogging.ShowStatus,
			IP:     config.RequestLogging.ShowIP,
			Method: config.RequestLogging.ShowMethod,
			Path:   config.RequestLogging.ShowPath,
			Query:  config.RequestLogging.ShowQuery,

			MessageContextKeys: []string{"logger_message"},
			MessageHeaderKeys:  []string{"User-Agent"},
		})

		app.Use(requestLogger)
	}

	app.OnErrorCode(404, FourOhFour)

	if config.UseSSL {
		err := app.Run(iris.AutoTLS(":443", config.Domain, config.Email))

		if err != nil {
			panic(err)
		}

	} else {
		err := app.Run(iris.Addr(fmt.Sprintf(":%d", config.Port)))

		if err != nil {
			panic(err)
		}

	}

}

func FourOhFour(ctx iris.Context) {
	ctx.StatusCode(http.StatusNotFound)
	_, _ = ctx.Text("404 Not Found")
}
