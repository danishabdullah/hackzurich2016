package getrational

import (
	"encoding/json"
	"html/template"
	"log"
	"math/rand"
	"net/http"
)

func indexHandler(templ *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templ.ExecuteTemplate(w, "index.html", nil)
	})
}

func playHandler(templ *template.Template, db QuestionDatabase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// show one question as demo
		idx := rand.Int() % len(db)
		selected := db[idx]
		log.Printf("Selected: %+v", selected)

		templ.ExecuteTemplate(w, "play.html", selected)
	})
}

func questionHandler(db QuestionDatabase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		selected := db.SelectRandom(10)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(selected); err != nil {
			log.Printf("Error writing JSON: %s", err)
		}
	})
}
