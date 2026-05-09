package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Donbassenok/go-lab/internal/handler"
	"github.com/Donbassenok/go-lab/internal/repository"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Помилка при читанні .env файлу: %v", err)
	}
	dbURL := viper.GetString("DB_URL")
	fmt.Printf("Підключаємося до БД за адресою: %s\n", dbURL)

	db, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		log.Fatalf("Не вдалося відкрити БД: %v", err)
	}
	defer db.Close()

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatalf("Не вдалося створити драйвер БД: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"sqlite3",
		driver,
	)
	if err != nil {
		log.Fatalf("Помилка ініціалізації міграцій: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Помилка при виконанні міграцій: %v", err)
	}
	fmt.Println("Міграції успішно виконані! База готова.")

	repo := repository.NewSQLitePlantRepo(db)

	plantHandler := handler.NewPlantHandler(repo)

	mux := http.NewServeMux()

	plantHandler.RegisterRoutes(mux)

	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	})

	port := ":8080"
	fmt.Printf("Сервер запускається на порту %s...\n", port)
	
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Помилка при запуску сервера: %v", err)
	}
}