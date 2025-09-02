package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func connectDB() *pgx.Conn {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erro carregando .env:", err)
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@localhost:5432/%s",
		os.Getenv("DBUSER"),
		os.Getenv("DBPASSWD"),
		os.Getenv("DBNAME"),
	)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatal("Erro de conexão:", err)
	}
	return conn
}

func connectdb() *sql.DB {
	if db != nil {
		return db // já conectado
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro carregando .env:", err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DBHOST"),
		os.Getenv("DBPORT"),
		os.Getenv("DBUSER"),
		os.Getenv("DBPASSWD"),
		os.Getenv("DBNAME"),
	)

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Conectado ao PostgreSQL!")
	return db
}
