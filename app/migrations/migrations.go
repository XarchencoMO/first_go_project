package migrations

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	//"log"
	"os"
	"path/filepath"
	"strings"
)

const migrationsDir = "database/migrations"

// Функция для создания таблицы миграций
func createMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			filename TEXT UNIQUE NOT NULL,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

// Функция для получения списка выполненных миграций
func getExecutedMigrations(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query("SELECT filename FROM migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	executed := make(map[string]bool)
	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			return nil, err
		}
		executed[filename] = true
	}
	return executed, rows.Err()
}

// Функция для выполнения миграции
func executeMigration(db *sql.DB, filename string) error {

	data, err := os.ReadFile(filepath.Join(migrationsDir, filename))

	if err != nil {
		return err
	}

	_, err = db.Exec(string(data))

	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO migrations (filename) VALUES ($1)", filename)

	return err
}

// TODO нужно причесать код и избавиться от дублирования connection-ов.
func RunMigrations() error {

	// Загрузка переменных из .env файла
	err := godotenv.Load()

	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}
	// Получение значений переменных окружения
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	connection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connection)

	if err != nil {
		return fmt.Errorf("error opening database connection: %v", err)
	}

	defer db.Close()
	// Создание таблицы миграций, если она не существует
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("error creating migrations table: %v", err)
	}
	// Получение списка выполненных миграций
	executedMigrations, err := getExecutedMigrations(db)

	if err != nil {
		return fmt.Errorf("error getting executed migrations: %v", err)
	}
	// Считывание всех миграционных файлов из директории
	files, err := os.ReadDir(migrationsDir)

	if err != nil {
		return fmt.Errorf("error reading migrations directory: %v", err)
	}

	for _, file := range files {

		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		if executedMigrations[file.Name()] {
			continue // Миграция уже выполнена
		}

		fmt.Printf("Applying migration: %s\n", file.Name())

		if err := executeMigration(db, file.Name()); err != nil {
			return fmt.Errorf("error applying migration %s: %v", file.Name(), err)
		}
	}

	return nil
}
