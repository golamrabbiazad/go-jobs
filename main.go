package main

import (
	"main/scrapper"
	"os"
	"strings"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/", handlerHome)
	e.POST("/scrape", handlerForm)
	e.Logger.Fatal(e.Start(":1323"))
}

func handlerHome(c echo.Context) error {
	return c.File("./template/home.html")
}

func handlerForm(c echo.Context) error {
	jobs := "jobs.csv"
	data := "job.csv"
	defer os.Remove(jobs)
	input := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
	scrapper.Scrape(input)
	return c.Attachment(jobs, data)
}
