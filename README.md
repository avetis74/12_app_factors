# 12-Factor Go Application

Это приложение следует принципам [12-факторной архитектуры](https://12factor.net/).

## Архитектура

### ✅ Реализованные факторы:

1. **Codebase** - Одна кодовая база в Git
2. **Dependencies** - Зависимости в `go.mod`
3. **Config** - Конфигурация через переменные окружения
4. **Backing services** - PostgreSQL и Redis как внешние сервисы
5. **Build, release, run** - Dockerfile и Makefile
6. **Processes** - Stateless процессы
7. **Port binding** - Привязка к порту через переменную окружения
8. **Concurrency** - Горутины и graceful shutdown
9. **Disposability** - Быстрый запуск и корректное завершение
10. **Dev/prod parity** - Docker для одинаковой среды
11. **Logs** - Структурированные логи в stdout
12. **Admin processes** - Система миграций с версионированием

## Запуск

### Локальная разработка

1. Скопируйте пример конфигурации:
```bash
cp env.example .env
```

2. Запустите с помощью Docker Compose:
```bash
make docker-run
```

### Переменные окружения

- `DATABASE_URL` - URL подключения к PostgreSQL
- `REDIS_URL` - URL подключения к Redis
- `SERVER_PORT` - Порт сервера (по умолчанию 8080)
- `ENV` - Окружение (development/staging/production)

## API Endpoints

- `GET /health` - Проверка состояния
- `GET /cache/stats` - Статистика Redis кеша
- `GET /users` - Получить всех пользователей (с кешированием)
- `POST /users` - Создать пользователя
- `GET /users/:id` - Получить пользователя по ID (с кешированием)
- `PUT /users/:id` - Обновить пользователя
- `DELETE /users/:id` - Удалить пользователя

## Команды разработки

```bash
make build       # Сборка
make run         # Запуск
make test        # Тестирование
make docker-run  # Запуск в Docker
make fmt         # Форматирование кода
```

## Миграции базы данных

```bash
# Применить все миграции
make migrate-up

# Откатить последнюю миграцию
make migrate-down

# Показать текущую версию
make migrate-status

# Создать новую миграцию
make migrate-create

# Полный сброс БД (осторожно!)
make migrate-reset
```

### Структура миграций

```
migrations/
├── 001_create_users_table.up.sql     # Создание таблицы users
├── 001_create_users_table.down.sql   # Откат создания таблицы
├── 002_add_user_status.up.sql        # Добавление поля status
└── 002_add_user_status.down.sql      # Откат добавления поля
```
