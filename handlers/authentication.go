package handlers

import (
	"NoobOJ/database"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("super-secret-key"))

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t := template.Must(template.ParseFiles(
			"templates/index.html",
			"templates/register.html",
			"templates/footer.html",
		))
		type PageData struct {
			Title     string
			Pusername string
			Logout    string
			Admin     bool
		}
		data := PageData{
			Title:     "Register - NoobOJ",
			Pusername: "Login",
			Logout:    "Register",
			Admin:     false,
		}
		t.ExecuteTemplate(w, "index.html", data)
		return
	}

	name := r.FormValue("name")
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirm := r.FormValue("confirm_password")

	if password != confirm {
		fmt.Fprint(w, "Passwords do not match")
		return
	}

	var exists int
	database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username=? OR email=?", username, email).Scan(&exists)
	if exists > 0 {
		fmt.Fprint(w, "Username or Email already exists")
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	_, err := database.DB.Exec("INSERT INTO users(name, username, email, password, user_type) VALUES (?, ?, ?, ?, 'user')",
		name, username, email, hash)
	if err != nil {
		fmt.Fprint(w, "Database error: "+err.Error())
		return
	}

	// http.Redirect(w, r, "/login", http.StatusSeeOther)
	session, _ := store.Get(r, "session")
	session.Values["username"] = username
	session.Values["user_type"] = "user"
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t := template.Must(template.ParseFiles(
			"templates/index.html",
			"templates/login.html",
			"templates/footer.html",
		))
		type PageData struct {
			Title     string
			Pusername string
			Logout    string
			Admin     bool
		}
		data := PageData{
			Title:     "Login - NoobOJ",
			Pusername: "Login",
			Logout:    "Register",
			Admin:     false,
		}
		t.ExecuteTemplate(w, "index.html", data)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	var username string
	var hash string
	var userType string

	err := database.DB.QueryRow("SELECT  username, password, user_type FROM users WHERE email=?", email).
		Scan(&username, &hash, &userType)
	if err == sql.ErrNoRows {
		fmt.Fprint(w, "Email not registered")
		return
	} else if err != nil {
		fmt.Fprint(w, "Database error: "+err.Error())
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Fprint(w, "Incorrect password")
		return
	}

	session, _ := store.Get(r, "session")
	session.Values["username"] = username
	session.Values["user_type"] = userType
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if session.Values["username"] == nil {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
