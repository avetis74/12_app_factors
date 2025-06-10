package handlers

import (
	"net/http"
	"strconv"

	"github.com/avetis74/12_app_factors/storage" // Импортируем наш пакет storage
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
		return c.JSON(http.StatusInternalServerError, "could not fetch users")
	}
	return c.JSON(http.StatusOK, users)
}

// CreateUser обрабатывает запрос на создание пользователя.
func (h *UserHandler) CreateUser(c echo.Context) error {
	var u storage.User
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid input")
	}

	if err := h.Store.CreateUser(&u); err != nil {
		// Здесь можно добавить более детальную обработку ошибок, например, дубликат email.
		return c.JSON(http.StatusInternalServerError, "could not create user")
	}

	return c.JSON(http.StatusCreated, u)
}

// GetUser, UpdateUser, DeleteUser реализуются по аналогии
func (h *UserHandler) GetUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid user id")
	}

	// Здесь нужно будет реализовать h.Store.GetUser(id)
	return c.JSON(http.StatusNotImplemented, fmt.Sprintf("get user with id %d not implemented", id))
}

