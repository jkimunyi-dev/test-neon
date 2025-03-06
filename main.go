package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Get the PostgreSQL connection URL from environment variables
	dbURL := os.Getenv("POSTGRES_URL")
	if dbURL == "" {
		log.Fatalf("POSTGRES_URL not set in .env file")
	}

	// Connect to the database
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	// Create a table (if not exists)
	_, err = conn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT);`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Insert a new user
	_, err = conn.Exec(context.Background(), `INSERT INTO users (name) VALUES ($1);`, "Alice")
	if err != nil {
		log.Fatalf("Failed to insert data: %v", err)
	}

	// Read and display users
	rows, err := conn.Query(context.Background(), `SELECT id, name FROM users;`)
	if err != nil {
		log.Fatalf("Failed to read data: %v", err)
	}
	defer rows.Close()

	fmt.Println("Users in DB:")
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		fmt.Printf("- ID: %d, Name: %s\n", id, name)
	}
}
