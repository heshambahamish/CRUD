package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"student-crud/db"
	"student-crud/models"

	"github.com/lib/pq"
)

func Index(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT * FROM students")
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var s models.Student
		if err := rows.Scan(&s.ID, &s.Name, &s.Email, &s.Age); err != nil {
			http.Error(w, "Error scanning data: "+err.Error(), http.StatusInternalServerError)
			return
		}
		students = append(students, s)
	}

	tmpl, err := template.ParseFiles("templates/index.html", "templates/layout.html")
	if err != nil {
		http.Error(w, "Template parsing error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", students)
}

func Create(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/create.html", "templates/layout.html")
	tmpl.ExecuteTemplate(w, "layout", nil)
}

func Store(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")
		ageStr := r.FormValue("age")

		age, err := strconv.Atoi(ageStr)
		if err != nil || age <= 0 {
			http.Error(w, "Invalid age", http.StatusBadRequest)
			return
		}

		if name == "" || email == "" {
			http.Error(w, "Name and Email are required", http.StatusBadRequest)
			return
		}

		_, err = db.DB.Exec("INSERT INTO students(name, email, age) VALUES($1, $2, $3)", name, email, age)
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Handle unique violation error for email
				tmpl, _ := template.ParseFiles("templates/create.html", "templates/layout.html")
				tmpl.ExecuteTemplate(w, "layout", map[string]interface{}{
					"Error": "A student with this email already exists",
				})
				return
			}
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	row := db.DB.QueryRow("SELECT * FROM students WHERE id=$1", id)
	var s models.Student
	row.Scan(&s.ID, &s.Name, &s.Email, &s.Age)

	tmpl, _ := template.ParseFiles("templates/edit.html", "templates/layout.html")
	tmpl.ExecuteTemplate(w, "layout", s)
}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id, _ := strconv.Atoi(r.FormValue("id"))
		name := r.FormValue("name")
		email := r.FormValue("email")
		ageStr := r.FormValue("age")

		age, err := strconv.Atoi(ageStr)
		if err != nil || age <= 0 {
			http.Error(w, "Invalid age", http.StatusBadRequest)
			return
		}

		if name == "" || email == "" {
			http.Error(w, "Name and Email are required", http.StatusBadRequest)
			return
		}

		_, err = db.DB.Exec("UPDATE students SET name=$1, email=$2, age=$3 WHERE id=$4", name, email, age, id)
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Handle unique violation error for email
				row := db.DB.QueryRow("SELECT * FROM students WHERE id=$1", id)
				var s models.Student
				row.Scan(&s.ID, &s.Name, &s.Email, &s.Age)

				tmpl, _ := template.ParseFiles("templates/edit.html", "templates/layout.html")
				tmpl.ExecuteTemplate(w, "layout", map[string]interface{}{
					"Error":   "A student with this email already exists",
					"Student": s,
				})
				return
			}
			http.Error(w, "Database update error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	db.DB.Exec("DELETE FROM students WHERE id=$1", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
