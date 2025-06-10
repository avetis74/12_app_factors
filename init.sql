-- Фактор XII: Admin processes - таблица для отслеживания миграций

-- Создаем таблицу для отслеживания версий миграций (golang-migrate создаст автоматически)
-- Но для примера покажем структуру:

-- CREATE TABLE IF NOT EXISTS schema_migrations (
--     version BIGINT NOT NULL PRIMARY KEY,
--     dirty BOOLEAN NOT NULL
-- );

-- Примечание: golang-migrate автоматически создает эту таблицу при первом запуске 