package main

import (
	"strconv"

	"github.com/elisaado/gluipertje-rewrite/config"
	"github.com/elisaado/gluipertje-rewrite/db"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := initRoutes()
	e.Use(middleware.CORS())

	config.InitConfig()

	db.InitDB()
	defer db.DB.Close()

	e.Logger.Fatal(e.Start(config.C.Host + ":" + strconv.Itoa(config.C.Port)))
}