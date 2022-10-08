package main

import (
	"database/sql"
	"flag"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

var db *sql.DB

func main() {
	flag.Parse()
	log.SetFlags(0)

	room := NewRoom()
	var err error
	db, err = sql.Open("sqlite3", "./src/db/database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", LoadHome)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/logout", Logout)
	http.HandleFunc("/admin", LoadAdmin)
	http.HandleFunc("/addtech", LoadAddTech)
	http.HandleFunc("/edittech", LoadEditTech)
	http.HandleFunc("/add", AddTech)
	http.HandleFunc("/edit", EditTech)
	http.HandleFunc("/delete", DelTech)
	http.HandleFunc("/tech", LoadTech)
	http.HandleFunc("/client", LoadClient)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(room, w, r)
	})
	http.HandleFunc("/404", Load404)

	fileServer := http.FileServer(http.Dir("./src/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	http.ListenAndServe(":8080", nil)
}
