package main

import (
	"html/template"
	"log"
	"net/http"
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

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.Handle("/play", playHandler(templ))
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
	context := struct {
		Question  string
		BoundLow  int
		BoundHigh int
		Unit      string
	}{
		Question:  "How many words are in the first Harry Potter Book?",
		BoundLow:  76944,
		BoundHigh: 76944,
		Unit:      "words",
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templ.ExecuteTemplate(w, "play.html", context)
	})
}
