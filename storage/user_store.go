package storage

import (
	"database/sql"
	"fmt"
)

// User описывает модель пользователя в базе данных.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserStore определяет интерфейс для работы с хранилищем пользователей.
// Использование интерфейса позволяет нам легко подменять реализацию (например, для тестов).
type UserStore interface {
	GetUsers() ([]User, error)
	GetUser(id int) (*User, error)
	CreateUser(user *User) error
	UpdateUser(id int, user *User) error
	DeleteUser(id int) error
}

// PostgresStore реализует интерфейс UserStore для работы с PostgreSQL.
type PostgresStore struct {
	DB *sql.DB
}

// NewPostgresStore создает новый экземпляр PostgresStore.
func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{DB: db}
}

// GetUsers возвращает всех пользователей из БД.
func (s *PostgresStore) GetUsers() ([]User, error) {
	rows, err := s.DB.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// CreateUser создает нового пользователя в БД.
// Возвращает ошибку, если что-то пошло не так. ID будет присвоен автоматически.
func (s *PostgresStore) CreateUser(user *User) error {
	// Мы возвращаем ID, чтобы обновить наш объект user
	err := s.DB.QueryRow(
		"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id",
		user.Name, user.Email,
	).Scan(&user.ID)

	return err
}

// GetUser, UpdateUser, DeleteUser реализуются аналогично...
// (Для краткости пока оставим их, вы можете добавить их по аналогии)

