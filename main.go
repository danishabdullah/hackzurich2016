package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

func safeHTML(text string) template.HTML {
	return template.HTML(text)
}

func main() {
	templ := template.New("root").Funcs(template.FuncMap{
		"safeHTML": safeHTML,
	})

	templ, err := templ.ParseGlob("templates/*")
	if err != nil {
		log.Fatalf("Can not parse templates: %s", err)
	}

	templ = templ.Funcs(template.FuncMap{
		"safeHTML": safeHTML,
	})

	questions, err := readDatabase()
	if err != nil {
		log.Fatalf("Can not read database: %s", err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.Handle("/play", playHandler(templ))
	http.Handle("/questions", questionHandler(questions))
	http.Handle("/", indexHandler(templ))

	addr := ":8080"
	log.Printf("Listening on %s...", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

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

// Question is the basic data entity.
type Question struct {
	Text      string `json:"text"`
	Unit      string `json:"unit"`
	BoundLow  uint64 `json:"boundLow"`
	BoundHigh uint64 `json:"boundHigh"`
}

func convertRecord(rec []string) (Question, error) {
	if len(rec) != 4 {
		return Question{}, fmt.Errorf("invalid length: %d", len(rec))
	}

	low, err := strconv.ParseFloat(rec[1], 64)
	if err != nil {
		return Question{}, err
	}

	high, err := strconv.ParseFloat(rec[2], 64)
	if err != nil {
		return Question{}, err
	}

	return Question{
		Text:      rec[0],
		Unit:      rec[3],
		BoundLow:  uint64(low),
		BoundHigh: uint64(high),
	}, nil
}

func readDatabase() ([]Question, error) {
	f, err := os.Open("Questions.csv")
	if err != nil {
		return []Question{}, err
	}

	reader := csv.NewReader(f)
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		return []Question{}, err
	}

	var result []Question
	for _, rec := range records {
		q, err := convertRecord(rec)
		if err != nil {
			log.Printf("Invalid record: %s", err)
			continue
		}

		result = append(result, q)
	}
	return result, nil
}

// QuestionDatabase is the interface for the database containing the questions.
type QuestionDatabase []Question

// SelectRandom selects `num` questions at random from the database.
func (db QuestionDatabase) SelectRandom(num int) []Question {
	if len(db) < num {
		return db
	}

	idx := rand.Perm(len(db))
	var result []Question
	for c, i := range idx {
		if c >= num {
			break
		}
		result = append(result, db[i])
	}

	return result
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
