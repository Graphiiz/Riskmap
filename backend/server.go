package main

import (
	d "backend/data"
	rmdb "backend/db"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	rmdb.DB()
	rmdb.Migrate()
	e := echo.New()

	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello")
	})

	e.GET("/SFLA", d.GetSFLAData)

	e.POST("/SFLA", d.WriteFromSFLA)

	e.POST("/MJM", d.UpdateFromMJM)

	e.POST("/RM/create", d.CreateRMData)

	e.GET("/RM", d.GetOverviewData)

	e.GET("/RM/filter_bar/:area", d.GetFilterBarData)

	e.Logger.Fatal(e.Start(":1234"))

}
