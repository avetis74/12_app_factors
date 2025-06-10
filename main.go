package main

import (
    "github.com/labstack/echo/v4"
    "github.com/YOURLOGIN/users/actions"
)

func main() {
    e := echo.New()

    // Routes
    e.GET("/users", actions.GetUsers)
    e.GET("/users/:id", actions.GetUser)
    e.POST("/users", actions.CreateUser)
    e.PUT("/users/:id", actions.UpdateUser)
    e.DELETE("/users/:id", actions.DeleteUser)

    // Start server
    e.Logger.Fatal(e.Start(":8080"))
}
