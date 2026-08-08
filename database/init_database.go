package database

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

var DB *sql.DB

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func InitDatabase() {
	var err error

	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASSWORD", "")
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "article")

	// 1. Hubungkan ke server MySQL tanpa nama database untuk memastikan database ada
	dbServerDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true", dbUser, dbPass, dbHost, dbPort)
	serverDB, err := sql.Open("mysql", dbServerDsn)
	if err != nil {
		log.Fatalf("Gagal membuka koneksi ke server MySQL: %v", err)
	}

	err = serverDB.Ping()
	if err != nil {
		log.Fatalf("Gagal terhubung ke server MySQL (ping): %v", err)
	}

	_, err = serverDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", dbName))
	if err != nil {
		log.Fatalf("Gagal membuat database %s: %v", dbName, err)
	}
	serverDB.Close()

	// 2. Hubungkan ke database target
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Gagal membuka koneksi ke database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Gagal terhubung ke database (ping): %v", err)
	}

	// 3. Jalankan migrasi database
	d, err := iofs.New(migrationFS, "migrations")
	if err != nil {
		log.Fatalf("Gagal inisialisasi driver iofs migrasi: %v", err)
	}

	driver, err := mysql.WithInstance(DB, &mysql.Config{})
	if err != nil {
		log.Fatalf("Gagal inisialisasi driver mysql migrasi: %v", err)
	}

	m, err := migrate.NewWithInstance("iofs", d, dbName, driver)
	if err != nil {
		log.Fatalf("Gagal membuat instance migrasi: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Gagal menjalankan migrasi: %v", err)
	}

	log.Println("Database dan migrasi tabel posts berhasil diinisialisasi!")
}
