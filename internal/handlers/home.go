package handlers

import (
	"html/template"
	"log"
	"net/http"
)

// ServeHome serves the home page
func ServeHome(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving home page")
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("executing template")
	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Done serving home page")
}
