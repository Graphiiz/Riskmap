package main

import (
	d "backend/data"
	rmdb "backend/db"
	bearer "backend/middleware"
	auth "backend/services"
	"net/http"
	"os"

	// "time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	// "github.com/robfig/cron"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title RiskMap API
// @version 0.1
// @description RiskMap Server for PEA.

// @BasePath /
// @schemes http
func main() {
	rmdb.DB()
	rmdb.Migrate()
	e := echo.New()

	e.Use(middleware.CORS())
	e.Pre(bearer.BearerToken)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "None Action")
	})

	e.GET("/auth", auth.GetAuthorizationUrl)

	e.GET("/auth/logout", auth.GetLogout)

	e.POST("/auth/token", auth.GetTokenClient)

	e.GET("/auth/user/info", auth.GetUserInfo)

	e.GET("/auth/image/user", auth.GetImage)

	sfla := e.Group("SFLA")

	sfla.GET("/SFLA", d.GetSFLAData)

	sfla.POST("/SFLA", d.WriteFromSFLA)

	e.POST("/MJM", d.UpdateFromMJM)

	e.POST("/RM/create", d.CreateRMData)

	e.GET("/RM", d.GetOverviewData)

	e.GET("/RM/filter_bar/:area", d.GetFilterBarData)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// c := cron.New()
	// c.AddFunc("* * * * * *", RunJob)
	// c.AddFunc("01 00 * * 1", RunJob)
	// go c.Start()

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))

}

// func RunJob() {

// 	fmt.Println("cron job run ......")
// 	fmt.Printf("%v\n", time.Now())

// 	// RunJobByShellScript()
// 	// RunJobByPython()
// 	// RunJobByPython()
// }

// func RunJobByPythonSayHello() {
// 	cmd := exec.Command("python", "-c", "import hello;")
// 	fmt.Println(cmd.Args)
// 	out, err := cmd.CombinedOutput()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(string(out))
// }

// func RunJobByShellScript() {
// 	cmd := exec.Command("/bin/sh", "./First.sh")
// 	cmd.Stdin = strings.NewReader("")
// 	var out bytes.Buffer
// 	cmd.Stdout = &out
// 	err := cmd.Run()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Output", out.String())
// }

// func RunJobByPython() {
// 	cmd := exec.Command("python", "-c", "import pythonFile; pythonFile.cat_strings('foo', 'bar')")
// 	fmt.Println(cmd.Args)
// 	out, err := cmd.CombinedOutput()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(string(out))
// }
