package main

import (
	"crypto/sha256"
	"database/sql"
	"log"
)

func userExists(db *sql.DB, username string) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}
	return exists
}

func userIsAdmin(db *sql.DB, username string) bool {
	var isAdmin bool
	err := db.QueryRow("SELECT is_admin FROM users WHERE username = ?", username).Scan(&isAdmin)
	if err != nil {
		log.Fatal(err)
	}
	return isAdmin
}

func addUser(db *sql.DB, username string, password string, isAdmin int) {

	hashedPassword := hashPassword(password)

	query := "INSERT INTO users (username, password, is_admin) VALUES (?, ?, ?)"

	db.Prepare(query)

	_, err := db.Exec(query, username, hashedPassword, isAdmin)
	if err != nil {
		log.Fatal(err)
	}
}

func removeUser(db *sql.DB, username string) {
	query := "DELETE FROM users WHERE username = ?"
	db.Prepare(query)
	_, err := db.Exec(query, username)
	if err != nil {
		log.Fatal(err)
	}
}

func checkPassword(db *sql.DB, username string, password string) bool {
	var hashedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
	if err != nil {
		log.Fatal(err)
	}
	return hashedPassword == hashPassword(password)
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return string(hash[:])
}
