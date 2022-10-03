package main

import (
	"database/sql"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var upgrader = websocket.Upgrader{}
var db *sql.DB

func main() {

	db, err := sql.Open("sqlite3", "./src/db/database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", loadHome)
	http.HandleFunc("/client", loadClient)
	http.HandleFunc("/404", loadClient)
	http.HandleFunc("/ws", ws)

	fileServer := http.FileServer(http.Dir("./src/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	http.ListenAndServe(":8080", nil)

}

func loadHome(w http.ResponseWriter, r *http.Request) {

	vue, _ := os.ReadFile("./src/views/template.html")
	vueStr := string(vue)
	home, _ := os.ReadFile("./src/views/home.html")
	homeStr := string(home)

	if r.URL.Path != "/" {
		homeStr = strings.Replace(homeStr, "###ERROR###", "d-block", 1)
		homeStr = strings.Replace(homeStr, "###MSG###", "Assurez-vous que votre nom d'utilisateur et votre mot de passe sont valides !", 1)
	} else {
		homeStr = strings.Replace(homeStr, "###ERROR###", "d-none", 1)
	}

	vueStr = strings.Replace(vueStr, "###TITLE###", "Clavardage du C.A.I.", 1)
	vueStr = strings.Replace(vueStr, "###SUBTITLE###", "On vous aide même à distance !", 1)
	vueStr = strings.Replace(vueStr, "###CONTENT###", homeStr, 1)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, vueStr)
}

func loadClient(w http.ResponseWriter, r *http.Request) {

	vue, _ := os.ReadFile("./src/views/template.html")
	vueStr := string(vue)
	client, _ := os.ReadFile("./src/views/client.html")
	clientStr := string(client)
	vueStr = strings.Replace(vueStr, "###TITLE###", "Clavardage du C.A.I.", 1)
	vueStr = strings.Replace(vueStr, "###SUBTITLE###", "Attendez, un technicien est sur le point de vous aider !", 1)
	vueStr = strings.Replace(vueStr, "###CONTENT###", clientStr, 1)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, vueStr)
}

func loadTech(w http.ResponseWriter, r *http.Request) {

	vue, _ := os.ReadFile("./src/views/template.html")
	vueStr := string(vue)
	home, _ := os.ReadFile("./src/views/tech.html")
	homeStr := string(home)
	vueStr = strings.Replace(vueStr, "###TITLE###", "Clavardage du C.A.I.", 1)
	vueStr = strings.Replace(vueStr, "###SUBTITLE###", "Aidez le plus de personne possible !", 1)
	vueStr = strings.Replace(vueStr, "###CONTENT###", homeStr, 1)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, vueStr)
}

func load404(w http.ResponseWriter, r *http.Request) {

	vue, _ := os.ReadFile("./src/views/template.html")
	vueStr := string(vue)
	home, _ := os.ReadFile("./src/views/404.html")
	homeStr := string(home)
	vueStr = strings.Replace(vueStr, "###TITLE###", "Clavardage du C.A.I.", 1)
	vueStr = strings.Replace(vueStr, "###SUBTITLE###", "Page introuvable", 1)
	vueStr = strings.Replace(vueStr, "###CONTENT###", homeStr, 1)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, vueStr)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		load404(w, r)
		return
	}
	switch r.Method {
	case "GET":
		load404(w, r)
	case "POST":
		user := r.FormValue("username")
		pass := r.FormValue("password")
		if userExists(db, user) && checkPassword(db, user, pass) {
			loadTech(w, r)
		} else {
			loadHome(w, r)
			io.WriteString(w, "Login failed")
		}
	default:
		load404(w, r)
	}
}

func ws(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
	}
	defer conn.Close()
	for {
		mt, msg, _ := conn.ReadMessage()
		log.Printf("recv: %s", msg)
		conn.WriteMessage(mt, msg)
	}
}
