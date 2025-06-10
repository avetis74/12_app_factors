package storage

import (
	"database/sql"
	"fmt" // fmt теперь нужен, так как мы используем fmt.Errorf
)

// User описывает модель пользователя в базе данных.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserStore определяет интерфейс для работы с хранилищем пользователей.
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
func (s *PostgresStore) CreateUser(user *User) error {
	err := s.DB.QueryRow(
		"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id",
		user.Name, user.Email,
	).Scan(&user.ID)
	return err
}

// НОВЫЙ МЕТОД
// GetUser находит одного пользователя по ID.
func (s *PostgresStore) GetUser(id int) (*User, error) {
	var u User
	// QueryRow идеально подходит для запросов, которые должны вернуть не более одной строки.
	err := s.DB.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		// Важно проверять на sql.ErrNoRows, чтобы корректно обрабатывать случай,
		// когда пользователь просто не найден, а не когда произошла реальная ошибка.
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, err
	}
	return &u, nil
}

// НОВЫЙ МЕТОД
// UpdateUser обновляет данные пользователя по ID.
func (s *PostgresStore) UpdateUser(id int, user *User) error {
	// DB.Exec используется для команд, не возвращающих строки (UPDATE, DELETE, INSERT без RETURNING).
	result, err := s.DB.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3", user.Name, user.Email, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// Если ни одна строка не была затронута, значит пользователя с таким ID не существует.
	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}
	return nil
}

// НОВЫЙ МЕТОД
// DeleteUser удаляет пользователя по ID.
func (s *PostgresStore) DeleteUser(id int) error {
	result, err := s.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}
	return nil
}

