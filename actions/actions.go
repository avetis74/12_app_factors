package actions

import (
    "net/http"
    "github.com/labstack/echo/v4"
)

// Placeholder for user data
type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

var users = []User{
    {ID: "1", Name: "John Doe", Email: "john@example.com"},
    {ID: "2", Name: "Jane Doe", Email: "jane@example.com"},
}

// GetUsers returns a list of users
func GetUsers(c echo.Context) error {
    return c.JSON(http.StatusOK, users)
}

// GetUser returns a single user by ID
func GetUser(c echo.Context) error {
    id := c.Param("id")
    for _, user := range users {
        if user.ID == id {
            return c.JSON(http.StatusOK, user)
        }
    }
    return c.JSON(http.StatusNotFound, "User not found")
}

// CreateUser adds a new user
func CreateUser(c echo.Context) error {
    user := new(User)
    if err := c.Bind(user); err != nil {
        return c.JSON(http.StatusBadRequest, "Invalid input")
    }
    users = append(users, *user)
    return c.JSON(http.StatusCreated, user)
}

// UpdateUser updates an existing user
func UpdateUser(c echo.Context) error {
    id := c.Param("id")
    updatedUser := new(User)
    if err := c.Bind(updatedUser); err != nil {
        return c.JSON(http.StatusBadRequest, "Invalid input")
    }
    for i, user := range users {
        if user.ID == id {
            users[i] = *updatedUser
            return c.JSON(http.StatusOK, updatedUser)
        }
    }
    return c.JSON(http.StatusNotFound, "User not found")
}

// DeleteUser deletes a user by ID
func DeleteUser(c echo.Context) error {
    id := c.Param("id")
    for i, user := range users {
        if user.ID == id {
            users = append(users[:i], users[i+1:]...)
            return c.JSON(http.StatusOK, "User deleted")
        }
    }
    return c.JSON(http.StatusNotFound, "User not found")
}
