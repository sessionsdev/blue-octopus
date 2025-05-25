package profile

import (
	"html/template"
	"net/http"
)

// ServeHome serves the home page
func ServeHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ServePettijohnProfile(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/pettijohn.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
