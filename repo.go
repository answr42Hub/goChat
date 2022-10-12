package main

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
)

func UserExists(db *sql.DB, username string) bool {
	var exists string
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&exists)
	if err != nil {
		return false
	}
	return true
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

func IsOnline(db *sql.DB, username string) bool {
	var token string
	err := db.QueryRow("SELECT token FROM users WHERE username = ?", username).Scan(&token)
	if err != nil {
		return false
	}
	return token != ""
}

func TechIsOnline(db *sql.DB) bool {
	token := ""
	err := db.QueryRow("SELECT token FROM users WHERE is_admin = 0 AND token IS NOT NULL").Scan(&token)
	if err != nil {
		return false
	}
	return token != ""
}

func GetUser(db *sql.DB, token string) string {
	var user string
	err := db.QueryRow("SELECT username FROM users WHERE token = ?", token).Scan(&user)
	if err != nil {
		return ""
	}
	return user
}

func GetUserId(db *sql.DB, username string) string {
	var id string
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&id)
	if err != nil {
		return ""
	}
	return id
}

func GetUserByID(db *sql.DB, id string) string {
	var user string
	err := db.QueryRow("SELECT username FROM users WHERE id = ?", id).Scan(&user)
	if err != nil {
		return ""
	}
	return user
}

func EditUser(db *sql.DB, userID string, newname string, password string) {
	query := "UPDATE users SET username = ?, password = ? WHERE id = ?"
	db.Prepare(query)

	_, err := db.Exec(query, newname, HashPassword(password), userID)
	if err != nil {
		log.Fatal(err)
	}
}

func GetUsers(db *sql.DB) map[string]string {
	rows, err := db.Query("SELECT id, username FROM users WHERE is_admin = 0")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	users := make(map[string]string)

	for rows.Next() {
		var user string
		var id string
		err = rows.Scan(&id, &user)
		if err != nil {
			log.Fatal(err)
		}
		users[id] = user
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return users
}

func DelUser(db *sql.DB, id string) {
	query := "DELETE FROM users WHERE id = ?"
	db.Prepare(query)
	_, err := db.Exec(query, id)
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

	_, err := db.Exec(query, nil, username)
	if err != nil {
		log.Fatal(err)
	}
}

func RandStringBytes(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	bs := hash.Sum(nil)
	str := fmt.Sprintf("%x", bs)
	return str
}
