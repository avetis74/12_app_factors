# Фактор V: Build, release, run - многостадийная сборка
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарный файл
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Финальный образ
FROM alpine:latest

# Устанавливаем ca-certificates для HTTPS запросов
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарный файл
COPY --from=builder /app/main .

# Фактор VI: Processes - запуск как stateless процесс
EXPOSE 8080

CMD ["./main"] 