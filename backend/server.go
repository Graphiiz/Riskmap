package main

import (
	d "backend/data"
	rmdb "backend/db"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron"
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

	c := cron.New()
	c.AddFunc("5 * * * *", RunJob)
	go c.Start()

	e.Logger.Fatal(e.Start(":1234"))

}

func RunJob() {

	fmt.Println("cron job run ......")
	fmt.Printf("%v\n", time.Now())

	cmd := exec.Command("/bin/sh", "./First.sh")
	cmd.Stdin = strings.NewReader("")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Output", out.String())

}
