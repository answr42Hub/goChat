package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./src/db/database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	c := make(chan []byte)
	http.HandleFunc("/", LoadHome)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/logout", Logout)
	http.HandleFunc("/client", LoadClient)
	http.HandleFunc("/admin", LoadAdmin)
	http.HandleFunc("/addtech", LoadAddTech)
	http.HandleFunc("/edittech", LoadEditTech)
	http.HandleFunc("/add", AddTech)
	http.HandleFunc("/edit", EditTech)
	http.HandleFunc("/delete", DelTech)
	http.HandleFunc("/tech", LoadTech)
	http.HandleFunc("/404", Load404)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		Ws(w, r, c)
	})

	fileServer := http.FileServer(http.Dir("./src/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	http.ListenAndServe(":8080", nil)
}
