package main

import (
	"github.com/elisaado/gluipertje-rewrite/handlers"
	"github.com/labstack/echo"
)

func initRoutes() *echo.Echo {
	e := echo.New()

	api := e.Group("/api")
	api.GET("/users", handlers.GetUsers)
	api.POST("/users", handlers.CreateUser)
	api.GET("/user/:idOrUsername", handlers.GetUserByIdOrUsername)
	api.GET("/messages", handlers.GetMessages)
	api.POST("/messages", handlers.SendMessage)
	api.GET("/messages/:limit", handlers.GetMessagesByLimit)
	api.GET("/message/:id", handlers.GetMessageById)
	api.POST("/token", handlers.RevokeToken)
	api.GET("/:token/me", handlers.GetUserByToken)
	return e
}
