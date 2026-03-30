package handlers

import (
	"NoobOJ/database"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("super-secret-key"))

func NewProblemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles(
			"templates/index.html",
			"templates/new.html",
			"templates/footer.html",
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session, _ := store.Get(r, "session")
		username, _ := session.Values["username"].(string)
		type PageData struct {
			Title     string
			Pusername string
			Logout    string
			Admin     bool
		}
		data := PageData{
			Title:     "Create Problem - NoobOJ",
			Pusername: username,
			Logout:    "Logout",
			Admin:     true,
		}
		t.ExecuteTemplate(w, "index.html", data)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Parse error", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	statement := r.FormValue("statement")
	inputDesc := r.FormValue("input_desc")
	outputDesc := r.FormValue("output_desc")
	constraints := r.FormValue("constraints")
	tagsRaw := r.FormValue("tags")
	rating := r.FormValue("rating")

	// Start a Transaction
	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "Transaction error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Defer a rollback in case of any early returns
	defer tx.Rollback()
	session, _ := store.Get(r, "session")
	username, _ := session.Values["username"].(string)
	res, err := tx.Exec("INSERT INTO problems (title, statement, input, output, constraints, author) VALUES (?, ?, ?, ?, ?,?)", title, statement, inputDesc, outputDesc, constraints, username)
	if err != nil {
		fmt.Fprint(w, "Error inserting problem: "+err.Error())
		return
	}
	problemID, err := res.LastInsertId()
	if err != nil {
		fmt.Fprint(w, "Error getting last ID: "+err.Error())
		return
	}
	res, err = tx.Exec("INSERT INTO ratings(problem_id,rating) VALUES(?,?)", problemID, rating)
	if err != nil {
		fmt.Fprint(w, "Error inserting rating: "+err.Error())
		return
	}
	// Insert Tags (Splitting by comma)
	if tagsRaw != "" {
		tags := strings.Split(tagsRaw, ",")
		for _, tag := range tags {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				_, err = tx.Exec("INSERT INTO tags (problem_id, tag) VALUES (?, ?)", problemID, tag)
				if err != nil {
					fmt.Fprint(w, "Error inserting tags: "+err.Error())
					return
				}
			}
		}
	}

	testInputs := r.Form["test_input[]"]
	testOutputs := r.Form["test_output[]"]
	testTypes := r.Form["test_type[]"]

	for i := 0; i < len(testInputs); i++ {
		if testInputs[i] == "" && testOutputs[i] == "" {
			continue
		}

		_, err = tx.Exec("INSERT INTO test_cases (problem_id, input, output, type) VALUES (?, ?, ?, ?)", problemID, testInputs[i], testOutputs[i], testTypes[i])
		if err != nil {
			fmt.Fprint(w, "Error inserting test case: "+err.Error())
			return
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		fmt.Fprint(w, "Commit error: "+err.Error())
		return
	}

	http.Redirect(w, r, "/administration", http.StatusSeeOther)
}
