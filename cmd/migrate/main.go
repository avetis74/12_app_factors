package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	var (
		databaseURL = flag.String("database-url", "", "Database URL")
		migrationsPath = flag.String("path", "file://migrations", "Path to migrations")
		command = flag.String("command", "", "Command: up, down, force, version, create")
		steps = flag.Int("steps", 0, "Number of steps for up/down commands")
		version = flag.Int("version", 0, "Version for force command")
		name = flag.String("name", "", "Migration name for create command")
	)
	flag.Parse()

	// Получаем DATABASE_URL из переменной окружения, если не передана
	if *databaseURL == "" {
		*databaseURL = os.Getenv("DATABASE_URL")
	}

	if *databaseURL == "" {
		log.Fatal("Database URL is required. Set DATABASE_URL environment variable or use -database-url flag")
	}

	switch *command {
	case "create":
		if *name == "" {
			log.Fatal("Migration name is required for create command")
		}
		createMigration(*name)
	case "up", "down", "force", "version":
		runMigration(*databaseURL, *migrationsPath, *command, *steps, *version)
	default:
		fmt.Println("Available commands:")
		fmt.Println("  up [steps]     - Apply migrations")
		fmt.Println("  down [steps]   - Rollback migrations")
		fmt.Println("  force version  - Force database to version")
		fmt.Println("  version        - Show current version")
		fmt.Println("  create name    - Create new migration files")
		fmt.Println("\nExamples:")
		fmt.Println("  go run cmd/migrate/main.go -command=up")
		fmt.Println("  go run cmd/migrate/main.go -command=down -steps=1")
		fmt.Println("  go run cmd/migrate/main.go -command=version")
		fmt.Println("  go run cmd/migrate/main.go -command=create -name=add_user_role")
	}
}

func createMigration(name string) {
	// Найдем следующий номер миграции
	files, err := os.ReadDir("migrations")
	if err != nil {
		log.Fatalf("Error reading migrations directory: %v", err)
	}

	nextVersion := 1
	for _, file := range files {
		if !file.IsDir() {
			filename := file.Name()
			if len(filename) >= 3 {
				if version, err := strconv.Atoi(filename[:3]); err == nil {
					if version >= nextVersion {
						nextVersion = version + 1
					}
				}
			}
		}
	}

	// Создаем файлы миграции
	upFile := fmt.Sprintf("migrations/%03d_%s.up.sql", nextVersion, name)
	downFile := fmt.Sprintf("migrations/%03d_%s.down.sql", nextVersion, name)

	upContent := fmt.Sprintf(`-- Migration: %s
-- Version: %03d
-- Description: %s

-- TODO: Add your migration SQL here
`, name, nextVersion, name)

	downContent := fmt.Sprintf(`-- Rollback: %s
-- Version: %03d
-- Description: Rollback %s

-- TODO: Add your rollback SQL here
`, name, nextVersion, name)

	if err := os.WriteFile(upFile, []byte(upContent), 0644); err != nil {
		log.Fatalf("Error creating up migration: %v", err)
	}

	if err := os.WriteFile(downFile, []byte(downContent), 0644); err != nil {
		log.Fatalf("Error creating down migration: %v", err)
	}

	fmt.Printf("Created migration files:\n")
	fmt.Printf("  %s\n", upFile)
	fmt.Printf("  %s\n", downFile)
}

func runMigration(databaseURL, migrationsPath, command string, steps, version int) {
	// Открываем соединение с базой данных
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Создаем драйвер для postgres
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Error creating postgres driver: %v", err)
	}

	// Создаем мигратор
	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		log.Fatalf("Error creating migrator: %v", err)
	}

	// Выполняем команду
	switch command {
	case "up":
		if steps > 0 {
			err = m.Steps(steps)
		} else {
			err = m.Up()
		}
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Error applying migrations: %v", err)
		}
		if err == migrate.ErrNoChange {
			fmt.Println("No migrations to apply")
		} else {
			fmt.Println("Migrations applied successfully")
		}

	case "down":
		if steps > 0 {
			err = m.Steps(-steps)
		} else {
			err = m.Down()
		}
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Error rolling back migrations: %v", err)
		}
		if err == migrate.ErrNoChange {
			fmt.Println("No migrations to rollback")
		} else {
			fmt.Println("Migrations rolled back successfully")
		}

	case "force":
		if version == 0 {
			log.Fatal("Version is required for force command")
		}
		err = m.Force(version)
		if err != nil {
			log.Fatalf("Error forcing version: %v", err)
		}
		fmt.Printf("Forced database to version %d\n", version)

	case "version":
		currentVersion, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("Error getting version: %v", err)
		}
		fmt.Printf("Current version: %d\n", currentVersion)
		if dirty {
			fmt.Println("Database is in dirty state - migration failed")
		}
	}
} 