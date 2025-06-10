package storage

import (
	"database/sql"
	"fmt"
	"log"
)

// User описывает модель пользователя в базе данных.
type User struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
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
	log.Println("Fetching all users from database")
	rows, err := s.DB.Query("SELECT id, name, email, COALESCE(status, 'active') FROM users")
	if err != nil {
		log.Printf("Error querying users: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Status); err != nil {
			log.Printf("Error scanning user row: %v", err)
			return nil, err
		}
		users = append(users, u)
	}
	log.Printf("Successfully fetched %d users", len(users))
	return users, nil
}

// CreateUser создает нового пользователя в БД.
func (s *PostgresStore) CreateUser(user *User) error {
	log.Printf("Creating user: %s (%s)", user.Name, user.Email)
	
	// Устанавливаем статус по умолчанию, если не указан
	if user.Status == "" {
		user.Status = "active"
	}
	
	err := s.DB.QueryRow(
		"INSERT INTO users (name, email, status) VALUES ($1, $2, $3) RETURNING id",
		user.Name, user.Email, user.Status,
	).Scan(&user.ID)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}
	log.Printf("User created with ID: %d", user.ID)
	return nil
}

// GetUser находит одного пользователя по ID.
func (s *PostgresStore) GetUser(id int) (*User, error) {
	log.Printf("Fetching user with ID: %d", id)
	var u User
	err := s.DB.QueryRow("SELECT id, name, email, COALESCE(status, 'active') FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email, &u.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("User with ID %d not found", id)
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		log.Printf("Error fetching user %d: %v", id, err)
		return nil, err
	}
	log.Printf("Successfully fetched user: %s (%s)", u.Name, u.Email)
	return &u, nil
}

// UpdateUser обновляет данные пользователя по ID.
func (s *PostgresStore) UpdateUser(id int, user *User) error {
	log.Printf("Updating user %d: %s (%s)", id, user.Name, user.Email)
	
	// Устанавливаем статус по умолчанию, если не указан
	if user.Status == "" {
		user.Status = "active"
	}
	
	result, err := s.DB.Exec("UPDATE users SET name = $1, email = $2, status = $3 WHERE id = $4", user.Name, user.Email, user.Status, id)
	if err != nil {
		log.Printf("Error updating user %d: %v", id, err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected for user %d: %v", id, err)
		return err
	}
	if rowsAffected == 0 {
		log.Printf("User with ID %d not found for update", id)
		return fmt.Errorf("user with id %d not found", id)
	}
	log.Printf("Successfully updated user %d", id)
	return nil
}

// DeleteUser удаляет пользователя по ID.
func (s *PostgresStore) DeleteUser(id int) error {
	log.Printf("Deleting user with ID: %d", id)
	result, err := s.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		log.Printf("Error deleting user %d: %v", id, err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected for user deletion %d: %v", id, err)
		return err
	}
	if rowsAffected == 0 {
		log.Printf("User with ID %d not found for deletion", id)
		return fmt.Errorf("user with id %d not found", id)
	}
	log.Printf("Successfully deleted user %d", id)
	return nil
}

