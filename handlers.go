package getrational

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"

	"github.com/pborman/uuid"
)

// NumQuestions contains the number of questions for a single round.
const NumQuestions = 1

func initHandlers(mux *http.ServeMux, templ *template.Template, db QuestionDatabase) {
	mux.Handle("/api/questions/random", questionHandler(db))

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	mux.Handle("/play/", playHandler(templ, db))
	mux.Handle("/play", newGameHandler())
	mux.Handle("/game/", gameHandler(templ))
	mux.Handle("/", indexHandler(templ))
}

func indexHandler(templ *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templ.ExecuteTemplate(w, "index.html", nil)
	})
}

func newGameHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.NewRandom().String()
		// TODO: save game id somewhere

		http.Redirect(w, r, fmt.Sprintf("/play/%s", id), http.StatusFound)
	})
}

type playContext struct {
	ID        string
	Questions []Question
}

func playHandler(templ *template.Template, db QuestionDatabase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := path.Base(r.URL.Path)
		selected := db.SelectRandom(NumQuestions)

		templ.ExecuteTemplate(w, "play.html", playContext{
			ID:        id,
			Questions: selected,
		})
	})
}

func questionHandler(db QuestionDatabase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		selected := db.SelectRandom(NumQuestions)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(selected); err != nil {
			log.Printf("Error writing JSON: %s", err)
		}
	})
}

// Answer contains the information about an answer given by the user.
type Answer struct {
	Question   Question
	LowerBound uint64
	UpperBound uint64
}

type gameContext struct {
	Answers []Answer
}

func gameHandler(templ *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		defer r.Body.Close()

		var answers []Answer
		if err := json.NewDecoder(r.Body).Decode(&answers); err != nil {
			http.Error(w, fmt.Sprintf("Error parsing answers: %s", err), http.StatusBadRequest)
			return
		}

		templ.ExecuteTemplate(w, "game.html", gameContext{
			Answers: answers,
		})
	})
}
