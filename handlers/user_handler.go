package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/avetis74/12_app_factors/storage"
	"github.com/labstack/echo/v4"
)

// UserHandler содержит зависимости для обработчиков, в данном случае — хранилище.
type UserHandler struct {
	Store storage.UserStore
}

// NewUserHandler создает новый экземпляр UserHandler.
func NewUserHandler(s storage.UserStore) *UserHandler {
	return &UserHandler{Store: s}
}

// GetUsers обрабатывает запрос на получение всех пользователей.
func (h *UserHandler) GetUsers(c echo.Context) error {
	users, err := h.Store.GetUsers()
	if err != nil {
		// Фактор XI: Логи как потоки событий
		log.Printf("Error fetching users: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Could not fetch users",
		})
	}
	return c.JSON(http.StatusOK, users)
}

// CreateUser обрабатывает запрос на создание пользователя.
func (h *UserHandler) CreateUser(c echo.Context) error {
	var u storage.User
	if err := c.Bind(&u); err != nil {
		log.Printf("Error binding user data: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid input",
		})
	}

	if err := h.Store.CreateUser(&u); err != nil {
		log.Printf("Error creating user: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Could not create user",
		})
	}

	return c.JSON(http.StatusCreated, u)
}

// GetUser обрабатывает запрос на получение одного пользователя.
func (h *UserHandler) GetUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	user, err := h.Store.GetUser(id)
	if err != nil {
		log.Printf("Error fetching user %d: %v", id, err)
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	return c.JSON(http.StatusOK, user)
}

// UpdateUser обрабатывает запрос на обновление пользователя.
func (h *UserHandler) UpdateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	var u storage.User
	if err := c.Bind(&u); err != nil {
		log.Printf("Error binding user data: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid input",
		})
	}

	if err := h.Store.UpdateUser(id, &u); err != nil {
		log.Printf("Error updating user %d: %v", id, err)
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	u.ID = id
	return c.JSON(http.StatusOK, u)
}

// DeleteUser обрабатывает запрос на удаление пользователя.
func (h *UserHandler) DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	if err := h.Store.DeleteUser(id); err != nil {
		log.Printf("Error deleting user %d: %v", id, err)
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	return c.NoContent(http.StatusNoContent)
}