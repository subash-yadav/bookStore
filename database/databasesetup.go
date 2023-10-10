package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // postgres golang driver
)

func CreateConnection() *sql.DB {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("error loading .env file ")
	}

	// open the connection

	host := os.Getenv("DB_HOST")
	dbname := os.Getenv("DATABASE")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_HMAC")
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))

	psqlinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlinfo)

	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	// return the connection
	return db
}
