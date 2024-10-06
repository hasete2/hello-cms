package main

import (
	"hello-cms/db"
	"hello-cms/domain"
	"hello-cms/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.Use(middleware.Logger())

	db := db.NewDB()
	defer db.Close()

	c := domain.NewContentDomain(db)
	m := domain.NewManageDomain(db)
	h := handler.NewHandler(*c, *m)
	h.Register(e)

	e.Logger.Fatal(e.Start(":2345"))
}
