package getrational

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

func indexHandler(templ *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templ.ExecuteTemplate(w, "index.html", nil)
	})
}

func playHandler(templ *template.Template) http.Handler {
	context := Question{
		Text:      "How many words are in the first Harry Potter Book?",
		BoundLow:  76944,
		BoundHigh: 76944,
		Unit:      "words",
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templ.ExecuteTemplate(w, "play.html", context)
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
