package main

import (
	"crypto/sha256"
	"database/sql"
	"log"
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func UserExists(db *sql.DB, username string) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}
	return exists
}

func UserIsAdmin(db *sql.DB, username string) bool {
	var isAdmin bool
	err := db.QueryRow("SELECT is_admin FROM users WHERE username = ?", username).Scan(&isAdmin)
	if err != nil {
		log.Fatal(err)
	}
	return isAdmin
}

func AddUser(db *sql.DB, username string, password string, isAdmin int) {

	hashedPassword := HashPassword(password)

	query := "INSERT INTO users (username, password, is_admin) VALUES (?, ?, ?)"

	db.Prepare(query)

	_, err := db.Exec(query, username, hashedPassword, isAdmin)
	if err != nil {
		log.Fatal(err)
	}
}

func GetUser(db *sql.DB, token string) string {
	var user string
	err := db.QueryRow("SELECT username FROM users WHERE token = ?", token).Scan(&user)
	if err != nil {
		log.Fatal(err)
	}
	return user
}

func RemoveUser(db *sql.DB, username string) {
	query := "DELETE FROM users WHERE username = ?"
	db.Prepare(query)
	_, err := db.Exec(query, username)
	if err != nil {
		log.Fatal(err)
	}
}

func CheckPassword(db *sql.DB, username string, password string) bool {
	var hashedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
	if err != nil {
		log.Fatal(err)
	}
	return hashedPassword == HashPassword(password)
}

func Connect(db *sql.DB, username string, token string) {
	query := "UPDATE users SET token = ? WHERE username = ?"
	db.Prepare(query)

	_, err := db.Exec(query, token, username)
	if err != nil {
		log.Fatal(err)
	}
}

func Disconnect(db *sql.DB, username string) {
	query := "UPDATE users SET token = ? WHERE username = ?"
	db.Prepare(query)

	_, err := db.Exec(query, "NULL", username)
	if err != nil {
		log.Fatal(err)
	}
}

func IsConnected(db *sql.DB) bool {
	var connected bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE token IS NOT NULL)").Scan(&connected)
	if err != nil {
		log.Fatal(err)
	}
	return connected
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return string(hash[:])
}
