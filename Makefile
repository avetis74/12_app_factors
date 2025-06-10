# Фактор V: Build, release, run - автоматизация процессов

.PHONY: build run test clean docker-build docker-run dev

# Сборка приложения
build:
	go build -o bin/app .

# Запуск приложения
run: build
	./bin/app

# Тестирование
test:
	go test -v ./...

# Очистка
clean:
	rm -rf bin/

# Сборка Docker образа
docker-build:
	docker build -t 12-factor-app .

# Запуск через Docker Compose
docker-run:
	docker compose up --build

# Остановка Docker Compose
docker-stop:
	docker compose down

# Разработка с автоперезагрузкой (требует air: go install github.com/cosmtrek/air@latest)
dev:
	air

# Линтер
lint:
	golangci-lint run

# Форматирование кода
fmt:
	go fmt ./...

# Миграции
migrate-up:
	go run cmd/migrate/main.go -command=up

migrate-down:
	go run cmd/migrate/main.go -command=down -steps=1

migrate-status:
	go run cmd/migrate/main.go -command=version

migrate-create:
	@read -p "Enter migration name: " name; \
	go run cmd/migrate/main.go -command=create -name=$$name

# Сброс базы данных (осторожно!)
migrate-reset:
	go run cmd/migrate/main.go -command=down
	go run cmd/migrate/main.go -command=up 