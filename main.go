package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// Global variable for template and database
var tmpl *template.Template

func init() {

	// Parse all templates in folder templates
	tmpl, _ = template.ParseGlob("templates/*.html")

}

func main() {

	var err error
	// Connect to database
	db, err := sql.Open("mysql", "root:student@tcp(192.168.2.83:3306)/testdb")
	if err != nil {
		fmt.Println("Failed to connecto database", err)
		return
	}
	defer db.Close()

	http.HandleFunc("/", HomePage)
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		LoginPage(w, r, db)
	})
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		RegisterPage(w, r, db)
	})
	http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Failed to start http server", err)
	}

}

func RegisterPage(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var exists bool
		err := db.QueryRow("SELECT username FROM users WHERE username = ?", username).Scan(&exists)

		switch {

		case err == sql.ErrNoRows:
			stmt, err := db.Prepare("INSERT INTO users(username, password) VALUES(?, ?)")
			if err != nil {
				fmt.Println("Error preparing sql query", err)
			}
			defer stmt.Close()

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				fmt.Println("Failed to hash password", err)
			}

			_, err = stmt.Exec(username, hashedPassword)
			if err != nil {
				fmt.Println("Error adding data to database", err)
			}

			fmt.Println("Successfully added user")

		case err != sql.ErrNoRows:
			fmt.Println("User already exists", err)

		case err != nil:
			fmt.Println("Error writing to database, try again", err)

		}
	}

	// Show Regristration form
	err := tmpl.ExecuteTemplate(w, "register.html", nil)
	if err != nil {
		fmt.Println("Error loading register page", err)
	}

}

func LoginPage(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	if r.Method != "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var databaseUsername string
		var databasePassword string

		// Scan the input of username and find it in database
		err := db.QueryRow("SELECT username, password FROM users WHERE username=?", username).Scan(&databaseUsername, &databasePassword)

		if err != nil {
			fmt.Println("User already exists", err)
		}

		// Compare the hashed password from database to input password
		err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))

		if err != nil {
			fmt.Println("Password is incorrect", err)
			return
		}

		fmt.Println("Login successful")

	}

	// Show login form
	err := tmpl.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		fmt.Println("Failed to execute loginform", err)
	}

}

func HomePage(w http.ResponseWriter, r *http.Request) {
	// Render the login form template
	err := tmpl.ExecuteTemplate(w, "homePage.html", nil)
	if err != nil {
		fmt.Println("Failed to execute loginform", err)
	}

}
