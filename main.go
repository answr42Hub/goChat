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
	"time"
)

var upgrader = websocket.Upgrader{}
var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./src/db/database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", loadHome)
	http.HandleFunc("/login", login)
	http.HandleFunc("/client", loadClient)
	http.HandleFunc("/admin", loadTech)
	http.HandleFunc("/404", loadClient)
	http.HandleFunc("/ws", ws)

	//admin : techadminsuperpassword

	fileServer := http.FileServer(http.Dir("./src/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	http.ListenAndServe(":8080", nil)

}

func loadHome(w http.ResponseWriter, r *http.Request) {

	cookie, cookieError := r.Cookie("session")
	if cookieError == nil && IsConnected(db) {
		user := GetUser(db, cookie.Value)
		if user != "" {
			if UserIsAdmin(db, user) {
				http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)
				return
			} else {
				http.Redirect(w, r, "/tech", http.StatusTemporaryRedirect)
				return
			}
		}
	}
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
	cookie, cookieError := r.Cookie("session")
	usrMsg := "Bienvenue cher technicien "
	if cookieError == nil {
		usrMsg += GetUser(db, cookie.Value) + " !"
	}
	vueStr = strings.Replace(vueStr, "###SUBTITLE###", usrMsg, 1)
	vueStr = strings.Replace(vueStr, "###CONTENT###", homeStr, 1)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, vueStr)
}

func loadAdmin(w http.ResponseWriter, r *http.Request) {
	vue, _ := os.ReadFile("./src/views/template.html")
	vueStr := string(vue)
	admin, _ := os.ReadFile("./src/views/admin.html")
	adminStr := string(admin)
	vueStr = strings.Replace(vueStr, "###TITLE###", "Clavardage du C.A.I.", 1)
	vueStr = strings.Replace(vueStr, "###SUBTITLE###", "Bienvenue cher administrateur", 1)
	vueStr = strings.Replace(vueStr, "###CONTENT###", adminStr, 1)

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
		http.Redirect(w, r, "/404", http.StatusTemporaryRedirect)
		return
	}
	switch r.Method {
	case "GET":
		http.Redirect(w, r, "/404", http.StatusTemporaryRedirect)
	case "POST":
		user := r.FormValue("username")
		pass := r.FormValue("password")
		if UserExists(db, user) && CheckPassword(db, user, pass) {
			if UserIsAdmin(db, user) {
				token := HashPassword(RandStringBytes(32)) + HashPassword(user)
				Connect(db, user, token)
				sessionCookie := http.Cookie{Name: "session", Value: token, HttpOnly: true}
				http.SetCookie(w, &sessionCookie)
				http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)
			} else {
				token := HashPassword(RandStringBytes(32)) + HashPassword(user)
				Connect(db, user, token)
				sessionCookie := http.Cookie{Name: "session", Value: token, HttpOnly: true}
				http.SetCookie(w, &sessionCookie)
				http.Redirect(w, r, "/tech", http.StatusTemporaryRedirect)
			}
		} else {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}
	default:
		http.Redirect(w, r, "/404", http.StatusTemporaryRedirect)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	if IsConnected(db) {
		cookie, cookieError := r.Cookie("session")
		if cookieError == nil {
			user := GetUser(db, cookie.Value)
			if user != "" {
				Disconnect(db, user)
			}
			cookie.Value = "Unuse"
			cookie.Expires = time.Unix(0, 0)
			http.SetCookie(w, cookie)
		}
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
