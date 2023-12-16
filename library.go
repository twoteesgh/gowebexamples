package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var db *sql.DB

func main() {
	setup()

	// Register application routes
	r := mux.NewRouter()
	r.HandleFunc("/books/{author}/{title}", showBook).Methods("GET")

	http.ListenAndServe(":8008", r)
}

func setup() {
	// Dotenv initialization
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// Database initialization
	sqlUrl := fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?multiStatements=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	sqlDb, sqlErr := sql.Open("mysql", sqlUrl)
	if sqlErr != nil {
		panic(sqlErr)
	}
	if err := sqlDb.Ping(); err != nil {
		panic(err)
	}
	db = sqlDb

	// Run database migrations
	if driver, err := mysql.WithInstance(db, &mysql.Config{}); err != nil {
		panic(err)
	} else if m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	); err != nil {
		panic(err)
	} else {
		m.Up()
	}
}

func showBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "You are reading %s, by %s", vars["title"], vars["author"])
}