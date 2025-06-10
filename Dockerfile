# Фактор V: Build, release, run - многостадийная сборка
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Устанавливаем зависимости для сборки
RUN apk add --no-cache git

# Устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарный файл
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Собираем миграционный инструмент
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrate ./cmd/migrate

# Финальный образ
FROM alpine:latest

# Устанавливаем необходимые пакеты
RUN apk --no-cache add ca-certificates netcat-openbsd

WORKDIR /root/

# Копируем бинарные файлы и файлы миграций
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY --from=builder /app/migrations ./migrations
COPY docker-entrypoint.sh ./

# Делаем entrypoint исполняемым
RUN chmod +x docker-entrypoint.sh

# Фактор VI: Processes - запуск как stateless процесс
EXPOSE 8080

# Используем entrypoint для миграций
ENTRYPOINT ["./docker-entrypoint.sh"] 