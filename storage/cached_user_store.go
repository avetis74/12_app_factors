package storage

import (
	"context"
	"fmt"
	"log"
	"time"
)

// CachedUserStore обертка над UserStore с кешированием
type CachedUserStore struct {
	store UserStore
	cache CacheService
}

// NewCachedUserStore создает новый кешированный UserStore
func NewCachedUserStore(store UserStore, cache CacheService) *CachedUserStore {
	return &CachedUserStore{
		store: store,
		cache: cache,
	}
}

// GetUsers возвращает всех пользователей с кешированием
func (c *CachedUserStore) GetUsers() ([]User, error) {
	ctx := context.Background()
	cacheKey := "users:all"
	
	// Пытаемся получить из кеша
	var users []User
	err := c.cache.Get(ctx, cacheKey, &users)
	if err == nil {
		log.Println("Returning users from cache")
		return users, nil
	}

	// Если не найдено в кеше, получаем из БД
	log.Println("Cache miss, fetching users from database")
	users, err = c.store.GetUsers()
	if err != nil {
		return nil, err
	}

	// Сохраняем в кеш на 5 минут
	if cacheErr := c.cache.Set(ctx, cacheKey, users, 5*time.Minute); cacheErr != nil {
		log.Printf("Failed to cache users: %v", cacheErr)
	}

	return users, nil
}

// GetUser возвращает пользователя по ID с кешированием
func (c *CachedUserStore) GetUser(id int) (*User, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:%d", id)
	
	// Пытаемся получить из кеша
	var user User
	err := c.cache.Get(ctx, cacheKey, &user)
	if err == nil {
		log.Printf("Returning user %d from cache", id)
		return &user, nil
	}

	// Если не найдено в кеше, получаем из БД
	log.Printf("Cache miss, fetching user %d from database", id)
	userPtr, err := c.store.GetUser(id)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кеш на 10 минут
	if cacheErr := c.cache.Set(ctx, cacheKey, *userPtr, 10*time.Minute); cacheErr != nil {
		log.Printf("Failed to cache user %d: %v", id, cacheErr)
	}

	return userPtr, nil
}

// CreateUser создает пользователя и сбрасывает кеш
func (c *CachedUserStore) CreateUser(user *User) error {
	err := c.store.CreateUser(user)
	if err != nil {
		return err
	}

	ctx := context.Background()
	
	// Сбрасываем кеш списка всех пользователей
	if cacheErr := c.cache.Delete(ctx, "users:all"); cacheErr != nil {
		log.Printf("Failed to invalidate users cache: %v", cacheErr)
	}

	// Кешируем нового пользователя
	if user.ID > 0 {
		cacheKey := fmt.Sprintf("user:%d", user.ID)
		if cacheErr := c.cache.Set(ctx, cacheKey, *user, 10*time.Minute); cacheErr != nil {
			log.Printf("Failed to cache new user %d: %v", user.ID, cacheErr)
		}
	}

	return nil
}

// UpdateUser обновляет пользователя и сбрасывает кеш
func (c *CachedUserStore) UpdateUser(id int, user *User) error {
	err := c.store.UpdateUser(id, user)
	if err != nil {
		return err
	}

	ctx := context.Background()
	
	// Сбрасываем кеш конкретного пользователя
	cacheKey := fmt.Sprintf("user:%d", id)
	if cacheErr := c.cache.Delete(ctx, cacheKey); cacheErr != nil {
		log.Printf("Failed to invalidate user %d cache: %v", id, cacheErr)
	}

	// Сбрасываем кеш списка всех пользователей
	if cacheErr := c.cache.Delete(ctx, "users:all"); cacheErr != nil {
		log.Printf("Failed to invalidate users cache: %v", cacheErr)
	}

	return nil
}

// DeleteUser удаляет пользователя и сбрасывает кеш
func (c *CachedUserStore) DeleteUser(id int) error {
	err := c.store.DeleteUser(id)
	if err != nil {
		return err
	}

	ctx := context.Background()
	
	// Сбрасываем кеш конкретного пользователя
	cacheKey := fmt.Sprintf("user:%d", id)
	if cacheErr := c.cache.Delete(ctx, cacheKey); cacheErr != nil {
		log.Printf("Failed to invalidate user %d cache: %v", id, cacheErr)
	}

	// Сбрасываем кеш списка всех пользователей
	if cacheErr := c.cache.Delete(ctx, "users:all"); cacheErr != nil {
		log.Printf("Failed to invalidate users cache: %v", cacheErr)
	}

	return nil
} 