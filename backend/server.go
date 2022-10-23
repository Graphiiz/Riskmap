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

	// // for API connection test
	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello")
	// })

	// e.POST("/SFLA", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "SFLA API OK")
	// })

	// e.POST("/MJM", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "MJM API OK")
	// })

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello")
	})

	e.POST("/SFLA", d.UpdateFromSFLA)

	e.POST("/MJM", d.UpdateFromMJM)

	// for test migration, remove before going to production
	e.POST("/RM", d.UpdateRMData)

	e.GET("/RM/:area/:name/:status", d.GetOverviewData)

	e.Logger.Fatal(e.Start(":1234"))

}
