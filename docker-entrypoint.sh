#!/bin/sh

# Фактор XII: Admin processes - автоматическое применение миграций

set -e

echo "Starting migration process..."

# Ожидаем готовности базы данных
echo "Waiting for database to be ready..."
until nc -z ${DATABASE_HOST:-db} ${DATABASE_PORT:-5432}; do
  echo "Database is unavailable - sleeping"
  sleep 1
done

echo "Database is ready!"

# Применяем миграции
echo "Running database migrations..."
./migrate -command=up

if [ $? -eq 0 ]; then
    echo "Migrations completed successfully"
else
    echo "Migration failed!"
    exit 1
fi

# Запускаем основное приложение
echo "Starting application..."
exec ./main 